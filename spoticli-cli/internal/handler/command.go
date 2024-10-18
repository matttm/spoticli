package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/matttm/spoticli/spoticli-cli/internal/utilities"
)

func DownloadSong(title string) error {
	b, _ := utilities.GetBytesBackend(nil, "audio/proxy/1")
	// path := "~/Downloads/spoticli"
	// _ = os.MkdirAll(path, 0664)
	// filePath := fmt.Sprintf("%s/%s", path, "test.mp3")
	os.WriteFile("test.mp3", b, 0664)
	return nil
}
func StreamSong(title string) error {
	b, _ := utilities.GetBytesBackend(
		func(r *http.Request) {
			r.Header.Add("Range", "bytes=0-1000000")
		},
		"audio/proxy/stream/1",
	)
	fmt.Print(b)
	r := bytes.NewReader(b)
	closer := io.NopCloser(r)
	streamer, format, err := mp3.Decode(closer)
	if err != nil {
		panic(err)
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)

	shot := buffer.Streamer(0, buffer.Len())
	streamer.Close()
	speaker.Play(shot)
	select {}
	return nil
}
