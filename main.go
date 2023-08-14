package main

import (
	"fmt"

	"github.com/luqus/livespace/api"
)

func main() {

	// if err := godotenv.Load(".env"); err != nil {

	// }

	port := "8080"

	api := api.New()
	go api.VideoProcessorQueue.Run()
	api.Run(fmt.Sprintf(":%s", port))

}
