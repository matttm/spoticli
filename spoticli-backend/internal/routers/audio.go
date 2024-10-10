package routers

import (
	"github.com/gorilla/mux"
	"github.com/mstttm/spoticli/spoticli-backend/internal/controllers"
)

func AttachAudioRouter(audio *mux.Router) *mux.Router {

	get := audio.Path("/{id:[0-9]+}").Methods("GET")
	get.HandlerFunc(controllers.GetPresignedUrl)

	stream := audio.Path("/proxy/{id:[0-9]+}").Methods("GET")
	stream.HandlerFunc(controllers.GetAudio)

	// post := router.Methods("POST")
	// post.HandlerFunc(nil)

	return audio
}
