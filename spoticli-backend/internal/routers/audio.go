package routers

import (
	"github.com/gorilla/mux"
	"github.com/matttm/spoticli/spoticli-backend/internal/controllers"
)

func AttachAudioRouter(audio *mux.Router) *mux.Router {

	get := audio.Path("/{id:[0-9]+}").Methods("GET")
	get.HandlerFunc(controllers.GetPresignedUrl)

	download := audio.Path("/proxy/{id:[0-9]+}").Methods("GET")
	download.HandlerFunc(controllers.GetAudio)

	stream := audio.Path("/proxy/stream/{id:[0-9]+}").Methods("GET")
	stream.HandlerFunc(controllers.GetAudioPart)

	upload := audio.Path("").Methods("POST")
	upload.HandlerFunc(controllers.UploadMusicThroughPresigned)
	// post := router.Methods("POST")
	// post.HandlerFunc(nil)

	return audio
}
