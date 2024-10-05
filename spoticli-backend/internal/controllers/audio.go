package controllers

import "net/http"

func GetAudio(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
}
