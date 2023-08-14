package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
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
	app := fiber.New(
		fiber.Config{
			BodyLimit: 100 * 1024 * 1024,
		},
	)

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
	api.app.Get("/stream/:filename", api.ViewVideo)
	api.app.Get("/metadata", api.fetchMetaData)

	// TODO: authorization middleware
	api.app.Use("/api", middleware.Authorization)

	// TODO: upload video request
	api.app.Post("/api/upload/:videoID", api.handleChunks)
	api.app.Post("/api/metadata", api.createMetaData)
	api.app.Post("/api/upload/process/:videoID", api.processVideo)

	return api.app.Listen(addr)
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

func (api *Api) createMetaData(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	log.Println("create meta data api")

	// TODO: get uid from Authorization
	uid, ok := ctx.Locals("uid").(string)
	if !ok {
		print("uid interface to string failed")
	}
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
	creator, err := api.authentication.AuthenticationStore().FetchUserByUID(c, uid)
	if err != nil {
		log.Print(err)
		return ctx.Status(http.StatusInternalServerError).JSON("failed to fetch creator")
	}

	uniqueId := storage.NewID()
	metaData.ID = uniqueId.ID()
	metaData.VideoID = uniqueId.String()
	metaData.Creator = creator.FormatResponse()
	metaData.VideoUrl = fmt.Sprintf("http://127.0.0.1:8080/media/video/%s", metaData.VideoID)

	err = api.metaDataStore.CommitMetaData(c, metaData)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON("we are fucked")
	}

	ctx.Set("Access-Control-Allow-Origin", "*")

	return ctx.Status(http.StatusCreated).JSON(metaData.VideoID)
}

func (api *Api) handleChunks(ctx *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("handle chunks request")

	// TODO: get video id from request params
	videoID := ctx.Params("videoID")
	if videoID == "" {
		print("no video id provided")
		return ctx.Status(http.StatusBadRequest).JSON("invalid video id")
	}

	// TODO: get multipart file from request
	file, err := ctx.FormFile("file")
	if err != nil {
		fmt.Println(err)
		ctx.Status(http.StatusBadRequest)
		return ctx.Send([]byte(err.Error()))
	}

	// TODO: save the file to the file system
	filename := fmt.Sprintf("temp/videos/%s.mp4", videoID)
	fileData, _ := file.Open()
	b := make([]byte, 15*1024*1024)
	fileData.Read(b)

	if err := os.WriteFile(filename, b, 0644); err != nil {
		// An error occurred.
		ctx.Status(http.StatusInternalServerError)
		return ctx.Send([]byte(err.Error()))

	}
	ctx.Set("Access-Control-Allow-Origin", "*")

	return ctx.Status(http.StatusCreated).JSON("File uploaded successfully!")
}

func (api *Api) processVideo(ctx *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	videoID := ctx.Params("videoID")
	if videoID == "" {
		return ctx.Status(http.StatusBadRequest).JSON("invalid video id")
	}

	api.VideoProcessorQueue.AddProcess(videoID)
	ctx.Set("Access-Control-Allow-Origin", "*")

	return ctx.Status(http.StatusOK).JSON("upload complete")
}

func (api *Api) ViewVideo(ctx *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filename := ctx.Params("filename")
	if filename == "" {
		return ctx.Status(http.StatusBadRequest).JSON("invalid video id")
	}
	filePath := fmt.Sprintf("%s/%s/%s", mediaPath, "64d62218f8e3ee40e1feea48", filename)

	ctx.Set("Content-Type", "text/plain;charset=utf-8")
	ctx.Set("Access-Control-Allow-Origin", "*")
	return ctx.SendFile(filePath)
}

// func (api *Api) uploadChunks(ctx *fiber.Ctx) error {
// 	_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
// 	defer cancel()

// 	log.Println("upload chunks request route")

// 	videoID := ctx.Params("videoID")
// 	if videoID == "" {
// 		return ctx.Status(http.StatusBadRequest).JSON("invalid video ID")
// 	}

// 	// TODO: check if videoID has  associated metaData
// 	// TODO: if not return error

// 	// TODO check if content-type = application-octet-stream
// 	contentType := ctx.Get("Content-Type")
// 	if contentType != "application-octet-stream" {
// 		return ctx.Status(http.StatusUnsupportedMediaType).JSON("unsupported media type provided")
// 	}

// 	body := ctx.Body()

// 	file, err := os.Create(fmt.Sprintf("temp/videos/%s.mp4", videoID))
// 	if err != nil {
// 		log.Println(err)
// 		return ctx.Status(http.StatusInternalServerError).JSON("failed to create video file")
// 	}

// 	defer file.Close()
// 	file.Write(body)

// 	return ctx.Status(http.StatusOK).JSON("video file successfully uploaded")

// }

func (api *Api) fetchMetaData(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	metas, err := api.metaDataStore.FetchVideos(c)
	if err != nil {
		log.Println(err)
		return ctx.Status(http.StatusInternalServerError).JSON("error fetching meta data")
	}

	return ctx.Status(http.StatusOK).JSON(metas)
}
