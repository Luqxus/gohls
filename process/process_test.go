package process

import (
	"testing"
)

var (
	processor VideoProcessorQueue = NewRedisVideoProcessorQueue()
)

func TestProcess(t *testing.T) {

	processor.Run()
	// processor.AddProcess("sabrina")

	// err := Process("sabrina")
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
