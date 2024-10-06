package routers

import (
	"github.com/gorilla/mux"
	"github.com/mstttm/spoticli/spoticli-backend/internal/controllers"
)

func AttachAudioRouter(audio *mux.Router) *mux.Router {

	get := audio.Path("/{song-id}").Methods("GET")
	get.HandlerFunc(controllers.GetPresignedUrl)

	stream := audio.Path("/{presigned-url}").Methods("GET")
	stream.HandlerFunc(controllers.GetAudio)

	// post := router.Methods("POST")
	// post.HandlerFunc(nil)

	return audio
}
