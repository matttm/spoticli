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
	w.Header().Add("Content-Length", strconv.FormatInt(*length, 10))
	w.Write(body)
}

func GetAudioPart(w http.ResponseWriter, r *http.Request) {
	println("streaming via proxy")
	id, _ := strconv.Atoi(
		mux.Vars(r)["id"],
	)
	rangeStr := r.Header["Range"][0]
	var start, end int64
	_, err := fmt.Sscan("bytes=%d-%d", rangeStr, start, end)
	// TODO: ensure end not gt file size
	body, length, err := services.GetAudioPart(id, rangeStr)
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "audio/mp3")
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
