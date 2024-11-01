// Package controllers is where the main parsing of headers
// is done and passed to the corresponding service
package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/matttm/spoticli/spoticli-backend/internal/services"
)

var audioService services.ApiAudioService = &services.AudioServiceWrap{}

// GetPresignedUrl gets a presigned url for an object in s3
// to be downloaded.
func GetPresignedUrl(w http.ResponseWriter, r *http.Request) {
	println("getting presigned url")
	id, _ := strconv.Atoi(
		mux.Vars(r)["id"],
	)
	url, err := audioService.GetPresignedUrl(id)
	if err != nil {
		panic(err)
	}
	w.Write([]byte(url))
}

// GetAudio gets the bytes of an entire object in s3
func GetAudio(w http.ResponseWriter, r *http.Request) {
	println("downloading via proxy")
	id, _ := strconv.Atoi(
		mux.Vars(r)["id"],
	)
	body, length, err := audioService.GetAudio(id)
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "audio/mp3")
	// NOTE: when you don't attach a content-length header, the server uses transfer-encoding chunked, which is a form of streaming
	w.Header().Add("Content-Length", strconv.FormatInt(*length, 10))
	w.Write(body)
}

// GetAudioPart gets a requested segmenyt of an object in s3
//
// The range being requested is to be specified by the \'Range\'
// header which will have a value like "bytes=0-1000000"
//
// # The response's body will contain a portion of bytes
//
// The response's body may or may not be the range specified
// and will be specified by the 'Content-Range' headerof the
// response, in the format od "bytes 0-10000/293872"
func GetAudioPart(w http.ResponseWriter, r *http.Request) {
	println("streaming via proxy")
	id, _ := strconv.Atoi(
		mux.Vars(r)["id"],
	)
	// NOTE: if there's no header, must be trying to get first bytes
	var start, end int64
	if len(r.Header["Range"]) == 0 {
		start = 0
		end = 0
	} else {
		rangeStr := r.Header["Range"][0]
		_, err := fmt.Sscanf(rangeStr, "bytes=%d-%d", &start, &end)
		if err != nil {
			panic(err)
		}
	}
	// TODO: ensure end not gt file size
	body, length, fileSize, err := audioService.StreamAudioSegment(id, &start, &end)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ContentLength %d\n", *length)
	fmt.Printf("FileSize %d\n", *fileSize)
	w.Header().Add("Content-Type", "audio/mp3")
	w.Header().Add("Content-Length", fmt.Sprintf("%d", *length))
	w.Header().Set(
		"Content-Range",
		fmt.Sprintf(
			"bytes %d-%d/%d",
			start,
			end,
			*fileSize,
		),
	)
	w.WriteHeader(http.StatusPartialContent)
	w.Write(body)
}
