package controllers

import (
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
