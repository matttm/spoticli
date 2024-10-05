package controllers

import "net/http"

func GetPresignedUrl(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
}
func GetAudio(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
}
