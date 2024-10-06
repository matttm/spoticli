package routers

import (
	"github.com/gorilla/mux"
	"github.com/mstttm/spoticli/spoticli-backend/internal/controllers"
)

func AudioRouter() *mux.Router {
	router := mux.
		NewRouter()
	audio := router.PathPrefix("/audio")

	get := audio.Path("/{song-id}").Methods("GET")
	get.HandlerFunc(controllers.GetPresignedUrl)

	stream := audio.Path("/{presigned-url}").Methods("GET")
	stream.HandlerFunc(controllers.GetAudio)

	// post := router.Methods("POST")
	// post.HandlerFunc(nil)

	return router
}
