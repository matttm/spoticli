package utilities

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/matttm/spoticli/spoticli-cli/internal/config"
	"github.com/matttm/spoticli/spoticli-cli/internal/models"
)

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
	_, err = fmt.Sscanf(contentRange, "bytes %d-%d/%d", &seg.StartByte, &seg.EndByte, &seg.TotalBytes)
	if err != nil {
		panic(err)
	}
	return res.Body
}
func GetBytesBackend(headerCb func(*http.Request), seg *models.AudioSegment, args ...interface{}) ([]byte, error) {
	body := GetBytesOpenReader(headerCb, seg, args...)
	data, err := io.ReadAll(body)
	if err != nil {
		panic(err)
	}
	defer body.Close()
	return data, nil
}
func GetBufferedAudioSegment(header []byte, id int, seg *models.AudioSegment) ([]byte, beep.StreamSeekCloser, beep.Format) {
	data, _ := GetBytesBackend(
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
	if len(header) != 0 {
		var copyHeader [128]byte
		for i := range header {
			copyHeader[i] = header[i]
		}
		data = append(copyHeader[:], data...)
	} else {
		header = data[:128] // rip off mp3 header + ID3v1 heder
	}
	reader := bytes.NewReader(data)
	dataStream := io.NopCloser(reader)
	streamer, format, err := mp3.Decode(dataStream)
	if err != nil {
		panic(err)
	}
	return header, streamer, format
}

func getClient() *http.Client {
	return new(http.Client)
}
