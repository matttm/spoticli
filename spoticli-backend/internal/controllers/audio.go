package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/mstttm/spoticli/spoticli-backend/internal/services"
)

func GetPresignedUrl(w http.ResponseWriter, r *http.Request) {
	println("getting presigned url")
	svc := services.GetStorageService()
	req, err := svc.GetPresignedUrl()
	if err != nil {
		panic(err)
	}
	b, _ := json.Marshal(req)
	w.Write(b)
}
func GetAudio(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
}
