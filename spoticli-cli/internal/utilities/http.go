package utilities

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/matttm/spoticli/spoticli-cli/internal/config"
	"github.com/matttm/spoticli/spoticli-cli/internal/models"
)

// GetBytesOpenReader gets an item from the backend and returns
// the responsre body as a fresh reader
//
// headerCb is a function taking an *http.Request that customies it
//
// seg an AudioSegment that will be updated with the response'
// "Content-Range" header, nost useful for partial data
//
// args is an arbitrary number of items to be used in string formatting
func GetBytesOpenReader(headerCb func(*http.Request), seg *models.AudioSegment, args ...interface{}) io.ReadCloser {
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
	_, err = fmt.Sscanf(contentRange, "bytes %d-%d/%d", &seg.StartByte, &seg.EndByte, &seg.SegmentLength)
	if err != nil {
		panic(err)
	}
	return res.Body
}

// GetBytesOpenReader gets an item from the backend and returns
// the response body as a []byte
//
// headerCb is a function taking an *http.Request that customies it
//
// seg an AudioSegment that will be updated with the response'
// "Content-Range" header, nost useful for partial data
//
// args is an arbitrary number of items to be used in string formatting
func GetBytesBackend(headerCb func(*http.Request), seg *models.AudioSegment, args ...interface{}) ([]byte, error) {
	body := GetBytesOpenReader(headerCb, seg, args...)
	data, err := io.ReadAll(body)
	if err != nil {
		panic(err)
	}
	defer body.Close()
	return data, nil
}

// GetBufferedAudioSegment gets a segment of queried song and
// returns the data as a streamer
//
// header
//
// id of audio track
//
// seg an AudioSegment that will be updated with the response'
// "Content-Range" header, nost useful for partial data

func GetBufferedAudioSegment(id string, seg *models.AudioSegment) (beep.StreamSeekCloser, beep.Format) {
	dataStream := GetBytesOpenReader(
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
		fmt.Sprintf("audio/proxy/stream/%s", id),
	)
	streamer, format, err := mp3.Decode(dataStream)
	if err != nil {
		panic(err)
	}
	return streamer, format
}

func getClient() *http.Client {
	return new(http.Client)
}
