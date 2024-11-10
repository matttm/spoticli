package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/matttm/spoticli/spoticli-backend/internal/services"
)

func GetAllFilesOfType(w http.ResponseWriter, r *http.Request) {
	cdStr := mux.Vars(r)["cd"]
	cd, _ := strconv.Atoi(cdStr)
	files := services.GetAllFilesOfType(cd)
	b, err := json.Marshal(files)
	if err != nil {
		panic(err)
	}
	w.Write(b)
}
