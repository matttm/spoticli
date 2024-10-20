package handler

import (
	"os"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
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
	streamer, format := utilities.GetBufferedAudioSegment(1, &seg) // this function call has a side affect on seg
	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	// then perform loop for remainder of song
	ticker := time.NewTicker(time.Second * 15)
	done := make(chan int)
	go func() {
		for {
			select {
			case <-ticker.C:
				// make seg point to next desired segment
				delta := seg.EndByte - seg.StartByte
				seg.StartByte = seg.EndByte + 1
				seg.EndByte = min(seg.StartByte+delta, seg.TotalBytes)
				streamer, _ = utilities.GetBufferedAudioSegment(1, &seg) // this function call has a side affect on seg
				buffer.Append(streamer)
				if seg.EndByte == seg.TotalBytes {
					return
				}

			case <-done:
				ticker.Stop()
				return
			}
		}
	}()

	shot := buffer.Streamer(0, buffer.Len())
	streamer.Close() // TODO: ENSURE ALL STREAMERS CLOSED
	speaker.Play(shot)
	select {}
	return nil
}
