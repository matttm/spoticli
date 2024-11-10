package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/matttm/spoticli/spoticli-backend/internal/services"
)

func GetAllFilesOfType(w http.ResponseWriter, r *http.Request) {
	files := services.GetAllFilesOfType(1)
	b, err := json.Marshal(files)
	if err != nil {
		panic(err)
	}
	w.Write(b)
}
