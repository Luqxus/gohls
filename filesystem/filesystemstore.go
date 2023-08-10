package filesystem

import (
	"fmt"
	"os"
)

type FileSystemStore interface {
	SetVideoFile(videoID string, videoData []byte) error
	RetrieveVideoPath(videoID string) string
}

type Ext4Store struct {
}

func NewExt4Store() *Ext4Store {
	return &Ext4Store{}
}

func (fs *Ext4Store) SetVideoFile(videoID string, videoData []byte) error {
	filePath := fmt.Sprintf("/temps/video/%s.mp4", videoID)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	// Write data to file
	_, err = file.Write(videoData)
	if err != nil {
		return err
	}

	return nil
}

func (fs *Ext4Store) RetrieveVideoPath(videoID string) string {
	// TODO: check if file exists
	return fmt.Sprintf("/temps/video/%s.mp4", videoID)
}
