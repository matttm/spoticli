package handler

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/manifoldco/promptui"
	"github.com/matttm/spoticli/spoticli-cli/internal/config"
	"github.com/matttm/spoticli/spoticli-cli/internal/models"
	"github.com/matttm/spoticli/spoticli-cli/internal/utilities"
)

func DownloadSong() error {
	_id, _ := promptSongs()
	id := strconv.Itoa(_id)
	fmt.Println(id)
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
func StreamSong() error {
	_id, _ := promptSongs()
	id := strconv.Itoa(_id)
	fmt.Println(id)
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/10))

	// A zero Queue is an empty Queue.
	done := make(chan bool)
	start := make(chan bool)
	var queue models.AudioSegmentQueue
	speaker.Play(&queue)

	// creating struct to follow boundaries
	seg := models.AudioSegment{StartByte: -1, EndByte: -1, SegmentLength: 0}
	var streamer beep.StreamSeekCloser
	// then perform loop for remainder of song
	ticker := time.NewTicker(time.Second * config.SECONDS_TO_WAIT_PER_FRAMES)
	go func() {
		for ; ; <-ticker.C {
			streamer, _ = utilities.GetBufferedAudioSegment(id, &seg) // this function call has a side affect on seg
			// start channel done
			fmt.Printf("Content-Range: %d-%d, %d\n", seg.StartByte, seg.EndByte, seg.SegmentLength)

			// And finally, we add the song to the queue.
			speaker.Lock()
			queue.Add(streamer)
			speaker.Unlock()
			if seg.StartByte == 0 {
				start <- true
				close(start)
			}

			if seg.EndByte == seg.SegmentLength {
				fmt.Println("Finished streaming song")
				// end channel
				if _, ok := <-start; !ok {
					queue.Add(beep.Callback(func() {
						done <- true
						close(done)
					}))
					return
				}
			}
			// make seg point to next desired segment
			delta := seg.EndByte - seg.StartByte
			seg.StartByte = seg.EndByte + 1
			seg.EndByte = min(seg.StartByte+delta, seg.SegmentLength+1)
		}
	}()
	_ = <-start
	fmt.Println("Song starts now")
	_ = <-done
	fmt.Println("Song finished playing")
	return nil
}

func UploadMusic(path string) error {
	fmt.Printf("opening %s\n", path)
	var wg sync.WaitGroup
	files := utilities.CollectFiles(path, ".mp3")
	for _, file := range files {
		wg.Add(1)
		go utilities.UploadFileViaPresign(file.Path, file.Size, &wg)
	}
	wg.Wait()
	fmt.Println("uploaded all songs")
	return nil
}

func ListFiles(cd int) error {
	files := utilities.GetAllFilesOfType(cd)
	for _, file := range files {
		fmt.Printf("%s\n", file.Key_name)
	}
	return nil
}
func promptSongs() (int, error) {
	cd := 1
	files := utilities.GetAllFilesOfType(cd)
	songNames := []string{}
	for _, f := range files {
		songNames = append(songNames, f.Key_name)
	}
	prompt := promptui.Select{
		Label: "Select Song",
		Items: songNames,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return -1, fmt.Errorf("Prompt failed %v\n", err)
	}
	for _, f := range files {
		if result == f.Key_name {
			return f.Id, nil
		}
	}
	return -1, fmt.Errorf("Key not found")
}
