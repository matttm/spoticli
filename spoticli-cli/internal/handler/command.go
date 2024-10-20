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
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/10))

	// A zero Queue is an empty Queue.
	var queue models.AudioSegmentQueue
	speaker.Play(&queue)

	// creating struct to follow boundaries
	seg := models.AudioSegment{StartByte: 0, EndByte: -1, TotalBytes: 0}
	var streamer beep.StreamSeekCloser
	var format beep.Format
	// then perform loop for remainder of song
	ticker := time.NewTicker(time.Second * 15)
	go func() {
		for ; ; <-ticker.C {
			if seg.EndByte >= seg.TotalBytes {
				return
			}
			streamer, format = utilities.GetBufferedAudioSegment(1, &seg) // this function call has a side affect on seg

			// The speaker's sample rate is fixed at 44100. Therefore, we need to
			// resample the file in case it's in a different sample rate.
			resampled := beep.Resample(4, format.SampleRate, sr, streamer)

			// And finally, we add the song to the queue.
			speaker.Lock()
			queue.Add(resampled)
			speaker.Unlock()
			// make seg point to next desired segment
			delta := seg.EndByte - seg.StartByte
			seg.StartByte = seg.EndByte + 1
			seg.EndByte = min(seg.StartByte+delta, seg.TotalBytes)

		}
	}()

	select {}
	ticker.Stop()
	return nil
}
