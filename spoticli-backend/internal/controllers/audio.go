package controllers

import (
	"fmt"
	"net/http"

	"github.com/mstttm/spoticli/spoticli-backend/internal/services"
)

func GetPresignedUrl(w http.ResponseWriter, r *http.Request) {
	println("getting presigned url")
	url, err := services.GetPresignedUrl("")
	if err != nil {
		panic(err)
	}
	w.Write([]byte(url))
}
func GetAudio(w http.ResponseWriter, r *http.Request) {
	println("downloading via proxy")
	body, err := services.GetAudio("bat_country.mp3")
	if err != nil {
		panic(err)
	}
	w.Write(body)
}

func GetAudioPart(w http.ResponseWriter, r *http.Request) {
	println("streaming via proxy")
	body, length, err := services.GetAudioPart("bat_country.mp3", "bytes=0-137439020734")
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "audio/mp3")
	w.Header().Set("Content-Range", fmt.Sprintf("bytes 0-137439020734/%d", length))
	w.WriteHeader(http.StatusPartialContent)
	w.Write(body)
}
