package handler

import (
	"fmt"
	"os"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/matttm/spoticli/spoticli-cli/internal/config"
	"github.com/matttm/spoticli/spoticli-cli/internal/models"
	"github.com/matttm/spoticli/spoticli-cli/internal/utilities"
)

func DownloadSong(id string) error {
	seg := models.AudioSegment{StartByte: 0, EndByte: 0, SegmentLength: 0}
	b, _ := utilities.GetBytesBackend(
		nil,
		&seg,
		fmt.Sprintf("audio/proxy/%s", id),
	)
	// path := "~/Downloads/spoticli"
	// _ = os.MkdirAll(path, 0664)
	// filePath := fmt.Sprintf("%s/%s", path, "test.mp3")
	os.WriteFile("test.mp3", b, 0664)
	return nil
}
func StreamSong(id string) error {
	fmt.Println(id)
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/10))

	// A zero Queue is an empty Queue.
	done := make(chan bool, 2)
	var queue models.AudioSegmentQueue
	speaker.Play(&queue)

	// creating struct to follow boundaries
	seg := models.AudioSegment{StartByte: 0, EndByte: -1, SegmentLength: 0}
	var streamer beep.StreamSeekCloser
	// then perform loop for remainder of song
	ticker := time.NewTicker(time.Second * config.SECONDS_TO_WAIT_PER_FRAMES)
	go func() {
		for ; ; <-ticker.C {
			streamer, _ = utilities.GetBufferedAudioSegment(id, &seg) // this function call has a side affect on seg
			// start channel done
			done <- true

			fmt.Printf("Content-Range: %d-%d, %d\n", seg.StartByte, seg.EndByte, seg.SegmentLength)

			// And finally, we add the song to the queue.
			speaker.Lock()
			queue.Add(streamer)
			speaker.Unlock()
			if seg.EndByte == seg.SegmentLength {
				fmt.Println("Finished streaming song")
				// end channel
				done <- true
				return
			}
			// make seg point to next desired segment
			delta := seg.EndByte - seg.StartByte
			seg.StartByte = seg.EndByte + 1
			seg.EndByte = min(seg.StartByte+delta, seg.SegmentLength+1)

		}
	}()

	<-done
	return nil
}
