package handler

import (
	"os"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/matttm/spoticli/spoticli-cli/internal/models"
	"github.com/matttm/spoticli/spoticli-cli/internal/utilities"
)

func DownloadSong(title string) error {
	seg := models.AudioSegment{StartByte: 0, EndByte: 0, TotalBytes: 0}
	b, _ := utilities.GetBytesBackend(nil, &seg, "audio/proxy/1")
	// path := "~/Downloads/spoticli"
	// _ = os.MkdirAll(path, 0664)
	// filePath := fmt.Sprintf("%s/%s", path, "test.mp3")
	os.WriteFile("test.mp3", b, 0664)
	return nil
}
func StreamSong(title string) error {
	// creating struct to follow boundaries
	seg := models.AudioSegment{StartByte: 0, EndByte: 0, TotalBytes: 0}
	streamer, format := utilities.GetBufferedAudioSegment(1, &seg)
	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)

	shot := buffer.Streamer(0, buffer.Len())
	streamer.Close()
	speaker.Play(shot)
	select {}
	return nil
}
