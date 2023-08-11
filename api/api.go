package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/luqus/livespace/authentication"
	"github.com/luqus/livespace/middleware"
	"github.com/luqus/livespace/process"
	"github.com/luqus/livespace/storage"
	"github.com/luqus/livespace/types"
)

var (
	mediaPath = "media/videos"
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

	// TODO: register user || create account
	api.app.Post("/register", api.authentication.RegisterUser)
	api.app.Post("/login", api.authentication.LoginUser)

	// TODO: view video
	api.app.Get("/stream/:videoID", api.ViewVideo)

	// TODO: authorization middleware
	api.app.Use(middleware.Authorization)

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
	// var metaData = new(types.MetaData)
	// err := ctx.BodyParser(metaData)
	// if err != nil {
	// 	return ctx.Status(http.StatusBadRequest).JSON("invalid request body")
	// }

	metaString := ctx.FormValue("metaData")
	if metaString == "" {
		return ctx.Status(http.StatusBadRequest).JSON("invalid video meta data")
	}

	metaData := new(types.MetaData)
	json.Unmarshal([]byte(metaString), metaData)

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

func (api *Api) FetchVideos(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// TODO: fetch videos from database
	videosMetaData, err := api.metaDataStore.FetchVideos(c)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON("error fetching videos")
	}

	// TODO: return videos to user
	return ctx.Status(http.StatusOK).JSON(videosMetaData)
}

func (api *Api) FetchVideoMetaData(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// TODO: get videoID from request params
	videoID := ctx.Params("videoID")
	if videoID == "" {
		return ctx.Status(http.StatusBadRequest).JSON("invalid video id")
	}

	// TODO: fetch video with id from database
	videoMetaData, err := api.metaDataStore.FetchVideoMetaData(c, videoID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON("error fetching video")
	}

	return ctx.Status(http.StatusFound).JSON(videoMetaData)
}

func (api *Api) ViewVideo(ctx *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	videoID := ctx.Params("videoID")
	if videoID == "" {
		return ctx.Status(http.StatusBadRequest).JSON("invalid video id")
	}
	filePath := fmt.Sprintf("%s/%s", mediaPath, videoID)

	ctx.Set("Content-Type", "text/plain;charset=utf-8")
	ctx.Set("Access-Control-Allow-Origin", "*")
	return ctx.SendFile(filePath)
}
