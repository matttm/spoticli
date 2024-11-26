package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/matttm/spoticli/spoticli-backend/internal/services"
)

var fmiService services.FileMetaInfoServiceApi = &services.FileMetaInfoServiceWrap{}

func GetAllFilesOfType(w http.ResponseWriter, r *http.Request) {
	cdStr := mux.Vars(r)["cd"]
	cd, _ := strconv.Atoi(cdStr)
	files := fmiService.GetAllFilesOfType(cd)
	b, err := json.Marshal(files)
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
