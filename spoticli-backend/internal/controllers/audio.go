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
	// NOTE: if there's no header, must be trying to get first bytes
	var start, end int
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
