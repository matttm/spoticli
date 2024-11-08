package services

type ApiAudioService interface {
	// GetPresignedUrl gets a presigned url
	// for downloading an object from s3
	GetPresignedUrl(id int) (string, error)
	GetAudio(id int) ([]byte, *int64, error)
	StreamAudioSegment(id int, start, end *int64) ([]byte, *int, *int64, error)
	UploadMusicThroughPresigned(resource string) string
}
