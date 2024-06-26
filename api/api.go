package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/luqus/livespace/authentication"
	"github.com/luqus/livespace/process"
	"github.com/luqus/livespace/storage"
	"github.com/luqus/livespace/types"
)

type Api struct {
	app                 *fiber.App
	VideoProcessorQueue process.VideoProcessorQueue
	metaDataStore       storage.MetaDataStore
	authentication      authentication.Authentication
}

func New() *Api {
	app := fiber.New()

	userAuthenticationStore, metaDataStore := storage.NewStore()
	return &Api{
		app:                 app,
		VideoProcessorQueue: process.NewRedisVideoProcessorQueue(),
		metaDataStore:       metaDataStore,
		authentication:      authentication.NewUserAuthentication(userAuthenticationStore),
	}
}

func (api *Api) Run(addr string) error {

	// go api.videoProcessorQueue.Run()

	// TODO: upload video request
	api.app.Post("/upload", api.uploadVideo)

	return api.app.Listen(addr)
}

// TODO: upload video route
func (api *Api) uploadVideo(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// TODO: get uid from Authorization
	uid := ctx.Locals("uid").(string)
	if uid == "" {
		return ctx.Status(http.StatusBadRequest).JSON("invalid user")
	}

	// TODO:get video metaData
	var metaData = new(types.MetaData)
	err := ctx.BodyParser(metaData)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON("invalid request body")
	}

	// TODO: fetch creator's user data
	creator, err := api.authentication.AuthenticationStore().FetchUser(c, uid)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON("user registration failed")
	}

	uniqueId := storage.NewID()
	metaData.ID = uniqueId.ID()
	// metaData.VideoID = uniqueId.String()
	metaData.Creator = creator.FormatResponse()
	metaData.VideoID = "sabrina"
	metaData.VideoUrl = fmt.Sprintf("http://127.0.0.1:3000/media/video/%s", metaData.VideoID)

	// get video
	// videoFile, err := ctx.FormFile("video")
	// if err != nil {
	// 	return ctx.Status(http.StatusBadRequest).JSON("error fetching video")
	// }

	// TODO: save video file to filesystem
	// err = ctx.SaveFile(videoFile, fmt.Sprintf("temps/videos/%s.mp4", metaData.VideoID))
	// if err != nil {
	// 	return ctx.Status(http.StatusInternalServerError).JSON("error uploading video")
	// }

	//TODO:  add video to processing queue
	err = api.VideoProcessorQueue.AddProcess(metaData.VideoID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON("error processing video")
	}

	// TODO: save meta data

	err = api.metaDataStore.CommitMetaData(c, metaData)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON("we are fucked")
	}

	defer cancel()
	return ctx.Status(http.StatusCreated).JSON("video upload successful")
}

// func (api *Api) FetchVideos(ctx *fiber.Ctx) error {
// 	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// }
