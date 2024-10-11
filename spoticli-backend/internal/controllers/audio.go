package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/mstttm/spoticli/spoticli-backend/internal/services"
)

func GetPresignedUrl(w http.ResponseWriter, r *http.Request) {
	println("getting presigned url")
	svc := services.GetStorageService()
	req, err := svc.GetPresignedUrl("bat_country.mp3")
	if err != nil {
		panic(err)
	}
	b, _ := json.Marshal(req)
	w.Write(b)
}
func GetAudio(w http.ResponseWriter, r *http.Request) {
	println("downloading via proxy")
	svc := services.GetStorageService()
	res, err := svc.DownloadFile("bat_country.mp3")
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		panic(err)
	}
	w.Write(body)
}
