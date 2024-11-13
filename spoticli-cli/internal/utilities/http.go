package utilities

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/matttm/spoticli/spoticli-cli/internal/config"
	"github.com/matttm/spoticli/spoticli-cli/internal/models"
	ext_models "github.com/matttm/spoticli/spoticli-models"
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

func UploadFileViaPresign(filepath string, wg *sync.WaitGroup) {
	fmt.Printf("Spawning thread to handle upload for %s\n", filepath)
	segs := strings.Split(filepath, "/")
	filename := segs[len(segs)-1]
	// TODO: move filename from url to body
	url := fmt.Sprintf("http://%s/audio/%s", config.SERVER_URL, filename)
	res, err := getClient().Post(url, "application/json", nil)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	url = string(b)
	fmt.Printf("Put to %s\n", url)
	// fmt.Printf("Opening %s\n", filepath)
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	stat, _ := file.Stat()
	//  fmt.Printf("file %s", file)
	//  defer file.Close() -- no needed as it should be being handled by 'Do' method
	req, err := http.NewRequest(http.MethodPut, url, file)
	if err != nil {
		panic(err)
	}
	req.Header["Content-Type"] = []string{"audio/mp3"}
	req.ContentLength = stat.Size()
	res, err = getClient().Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v object %v with presigned URL returned %v.", req.Method,
		filename, res.StatusCode)
	fmt.Println(strings.Repeat("-", 88))
	defer res.Body.Close()
	defer wg.Done()
}

func GetAllFilesOfType(cd int) []ext_models.FileMetaInfo {
	var ret []ext_models.FileMetaInfo
	// TODO: change to query param
	url := fmt.Sprintf("http://%s/files/%d", config.SERVER_URL, cd)
	res, err := getClient().Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &ret)
	if err != nil {
		panic(err)
	}
	return ret
}
func getClient() *http.Client {
	return new(http.Client)
}
