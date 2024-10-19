package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/matttm/spoticli/spoticli-backend/internal/services"
)

func GetPresignedUrl(w http.ResponseWriter, r *http.Request) {
	println("getting presigned url")
	id, _ := strconv.Atoi(
		mux.Vars(r)["id"],
	)
	url, err := services.GetPresignedUrl(id)
	if err != nil {
		panic(err)
	}
	w.Write([]byte(url))
}
func GetAudio(w http.ResponseWriter, r *http.Request) {
	println("downloading via proxy")
	id, _ := strconv.Atoi(
		mux.Vars(r)["id"],
	)
	body, length, err := services.GetAudio(id)
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "audio/mp3")
	// NOTE: when you don't attach a content-length header, the server uses transfer-encoding chunked, which is a form of streaming
	w.Header().Add("Content-Length", strconv.FormatInt(*length, 10))
	w.Write(body)
}

func GetAudioPart(w http.ResponseWriter, r *http.Request) {
	println("streaming via proxy")
	id, _ := strconv.Atoi(
		mux.Vars(r)["id"],
	)
	// TODO: ADD CHECK FOR IF RANGE NOT GIVEN (ASSUME 0?)
	rangeStr := r.Header["Range"][0]
	var start, end int
	_, err := fmt.Sscan("bytes=%d-%d", rangeStr, start, end)
	// TODO: ensure end not gt file size
	body, length, fileSize, err := services.StreamAudioSegment(id, &start, &end)
	if err != nil {
		panic(err)
	}
	fmt.Println(*length)
	fmt.Println(*fileSize)
	w.Header().Add("Content-Type", "audio/mp3")
	w.Header().Add("Content-Length", strconv.FormatInt(*length, 10))
	w.Header().Set(
		"Content-Range",
		fmt.Sprintf(
			"bytes %d-%d/%d",
			start,
			end,
			length,
		),
	)
	w.WriteHeader(http.StatusPartialContent)
	w.Write(body)
}
