package utilities

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/matttm/spoticli/spoticli-cli/internal/config"
	"github.com/matttm/spoticli/spoticli-cli/internal/models"
)

func GetBytesBackend(headerCb func(*http.Request), seg *models.AudioSegment, args ...interface{}) ([]byte, error) {
	var t interface{} = config.SERVER_URL
	slice := make([]interface{}, 1)
	slice[0] = t
	args = append(slice, args...)
	fmt.Println(fmt.Sprintf("http://%s/%s", args...))
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("http://%s/%s", args...),
		nil,
	)
	if err != nil {
		panic(err)
	}
	if headerCb != nil {
		headerCb(req)
	}
	res, err := getClient().Do(req)
	if err != nil {
		panic(err)
	}
	// TODO: REFACTOR AND VALIDATE
	contentRange := res.Header["Content-Range"][0]
	_, err = fmt.Sscanf(contentRange, "bytes %d-%d/%d", &seg.StartByte, &seg.EndByte, &seg.TotalBytes)
	if err != nil {
		panic(err)
	}
	var data []byte
	data, err = io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	return data, nil
}
func GetBufferedAudioSegment(id int, seg *models.AudioSegment) (beep.StreamSeekCloser, beep.Format) {
	b, _ := GetBytesBackend(
		func(r *http.Request) {
			// TODO: REMOVE HEADER ON FIRST SEND
			if seg.StartByte != 0 {
				r.Header.Add(
					"Range",
					fmt.Sprintf("bytes=%d-%d", seg.StartByte, seg.EndByte),
				)
			}
		},
		seg,
		fmt.Sprintf("audio/proxy/stream/%d", id),
	)
	fmt.Print(b)
	r := bytes.NewReader(b)
	closer := io.NopCloser(r)
	streamer, format, err := mp3.Decode(closer)
	if err != nil {
		panic(err)
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	return streamer, format
}

func getClient() *http.Client {
	return new(http.Client)
}
