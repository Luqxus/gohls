package authentication

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/luqus/livespace/helpers"
	"github.com/luqus/livespace/storage"
	"github.com/luqus/livespace/tokens"
	"github.com/luqus/livespace/types"
)

type Authentication interface {
	RegisterUser(ctx *fiber.Ctx) error
	LoginUser(ctx *fiber.Ctx) error
	AuthenticationStore() storage.AuthenticationStore
}

type UserAuthentication struct {
	authenticationStore storage.AuthenticationStore
}

func (auth *UserAuthentication) AuthenticationStore() storage.AuthenticationStore {
	return auth.authenticationStore
}

func NewUserAuthentication(authStore storage.AuthenticationStore) *UserAuthentication {
	return &UserAuthentication{
		authenticationStore: authStore,
	}
}

func (auth *UserAuthentication) RegisterUser(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// TODO: parse user from request body
	userinput := new(types.UserInput)
	err := ctx.BodyParser(userinput)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON("invalid user input")
	}

	// TODO: validate user input

	// TODO: check if email exists
	exists, err := auth.authenticationStore.CheckEmailExists(c, userinput.Email)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON("error checking email")
	}
	if exists {
		return ctx.Status(http.StatusBadRequest).JSON("email already exists")
	}

	// TODO: check if username exists
	exists, err = auth.authenticationStore.CheckUsernameExists(c, userinput.Username)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON("error checking username")
	}
	if exists {
		return ctx.Status(http.StatusBadRequest).JSON("username already exists")
	}

	user := new(types.User)
	userID := storage.NewID()
	user.ID = userID.ID()
	user.UID = userID.String()
	user.Email = userinput.Email
	user.Username = userinput.Username
	user.Password, _ = helpers.HashPassword(userinput.Password)

	// TODO: commit user data to user Database
	auth.authenticationStore.CreateUser(c, user)

	return ctx.Status(http.StatusCreated).JSON("user succesfully created")

}

func (auth *UserAuthentication) LoginUser(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	loginInput := new(types.LoginInput)

	// TODO: get login data from request body
	err := ctx.BodyParser(loginInput)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON("invalid login inputs")
	}

	// TODO: fetch user with email
	user, err := auth.authenticationStore.FetchUser(c, loginInput.Email)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON("wrong email or password")
	}

	// TODO: verify password
	err = helpers.VerifyPassword(user.Password, loginInput.Password)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON("wrong email or password")
	}

	// TODO: generate tokens
	signedToken, err := tokens.GenerateTokens(user.UID, user.Username)
	if err != nil {
		log.Panic(err)
	}

	ctx.Set("authorization", signedToken)

	// TODO: return user
	return ctx.Status(http.StatusFound).JSON(user.FormatResponse())
}
