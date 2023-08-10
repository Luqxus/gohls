package process

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/redis/go-redis/v9"
)

type VideoProcessorQueue interface {
	AddProcess(processID string) error
	Run()
}

type RedisVideoProcessorQueue struct {
	client    *redis.Client
	context   context.Context
	workers   int
	queueName string
}

// TODO: This function will initialize our processorQueue
func NewRedisVideoProcessorQueue() *RedisVideoProcessorQueue {
	redisAddr := "localhost:6379"
	redisPassword := ""
	return &RedisVideoProcessorQueue{
		queueName: "video_processing_queue",
		client: redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Username: "",
			Password: redisPassword,
			DB:       0,
		}),
		context: context.Background(),
		workers: 5,
	}
}

func (rq *RedisVideoProcessorQueue) Run() {

	log.Printf("redis client: %v", rq.client)
	for idx := 1; idx <= rq.workers; idx++ {
		go rq.worker(idx)
	}

	select {}
}

func (rq *RedisVideoProcessorQueue) AddProcess(processID string) error {
	log.Printf("process ID: %s added", processID)
	return rq.client.RPush(rq.context, rq.queueName, processID).Err()
}

func (rq *RedisVideoProcessorQueue) worker(workerID int) {
	log.Printf("worker : %d started", workerID)
	for {
		result, err := rq.client.BLPop(rq.context, 0, rq.queueName).Result()
		if err != nil {
			log.Printf("worker %d: failed to dequeue task %s", workerID, err.Error())
			continue
		}

		processID := result[1]

		err = Process(processID)
		if err != nil {
			log.Printf("video process at worker %d failed. processID: %s", workerID, processID)
			// TODO: alert Api of process result
		}

		// TODO: when process successful -> Delete temp video file
		os.Remove(fmt.Sprintf("temp/videos/%s.mp4", processID))
		fmt.Printf("worker %d : finished processing", workerID)
	}
}

func Process(filename string) error {

	// video to hls command
	cmd := exec.Command(
		"ffmpeg",
		"-i",
		"temp/videos/"+filename+".mp4",
		"-codec:",
		"copy",
		"-start_number",
		"0",
		"-hls_time",
		"10",
		"-hls_list_size",
		"0",
		"-f",
		"hls",
		"media/videos/"+filename+".m3u8",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// ffmpeg -i ${fileName} -codec: copy -start_number 0 -hls_time 10 -hls_list_size 0 -f hls ${dest}/${index}.m3u8
	// Execute cmd
	return cmd.Run()
}
