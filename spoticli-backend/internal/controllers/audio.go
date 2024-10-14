package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mstttm/spoticli/spoticli-backend/internal/services"
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
	body, err := services.GetAudio(id)
	if err != nil {
		panic(err)
	}
	w.Write(body)
}

func GetAudioPart(w http.ResponseWriter, r *http.Request) {
	println("streaming via proxy")
	id, _ := strconv.Atoi(
		mux.Vars(r)["id"],
	)
	body, length, err := services.GetAudioPart(id, "bytes=0-1000000")
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "audio/mp3")
	w.Header().Set("Content-Range", fmt.Sprintf("bytes 0-1000000/%d", length))
	w.WriteHeader(http.StatusPartialContent)
	w.Write(body)
}
