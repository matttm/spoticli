package handler

import (
	"os"

	"github.com/matttm/spoticli/spoticli-cli/internal/utilities"
)

func DownloadSong(title string) error {
	b, _ := utilities.GetBytesBackend("audio/proxy/1")
	// path := "~/Downloads/spoticli"
	// _ = os.MkdirAll(path, 0664)
	// filePath := fmt.Sprintf("%s/%s", path, "test.mp3")
	os.WriteFile("test.mp3", b, 0664)
	return nil
}
func StreamSong(title string) error {
	b, _ := utilities.GetBytesBackend("audio/proxy/stream/1")
	return nil
}
