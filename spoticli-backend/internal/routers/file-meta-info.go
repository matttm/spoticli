package routers

import (
	"github.com/gorilla/mux"
	"github.com/matttm/spoticli/spoticli-backend/internal/controllers"
)

func AttachFileMetaInfoRouter(files *mux.Router) *mux.Router {

	get := files.Path("/{id:[0-9]+}").Methods("GET")
	get.HandlerFunc(controllers.GetAllFilesOfType)

	return files
}
