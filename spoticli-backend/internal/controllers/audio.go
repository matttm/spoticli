package controllers

import (
	"net/http"

	"github.com/mstttm/spoticli/spoticli-backend/internal/services"
)

func GetPresignedUrl(w http.ResponseWriter, r *http.Request) {
	svc := services.GetStorageService()
}
func GetAudio(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
}
