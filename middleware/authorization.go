package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/luqus/livespace/tokens"
)

func Authorization(ctx *fiber.Ctx) error {
	// TODO: get token from request header
	signedToken := ctx.Get("authorization")
	if signedToken == "" {
		return ctx.Status(http.StatusUnauthorized).JSON("no authorization header provided")
	}

	// TODO: verify token
	uid, username, err := tokens.VerifyToken(signedToken)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON("invalid or expired authorization header")
	}

	// TODO: if valid add uid & username to context and serve resource
	ctx.Locals("uid", uid)
	ctx.Locals("username", username)
	return ctx.Next()
}
