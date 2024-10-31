package services

type AudioServiceWrap struct{}

// GetPresignedUrl gets a presigned url
// for downloading an object from s3
func (s *AudioServiceWrap) GetPresignedUrl(id int) (string, error)
func (s *AudioServiceWrap) GetAudio(id int) ([]byte, *int64, error)
func (s *AudioServiceWrap) StreamAudioSegment(id int, start, end *int64) ([]byte, *int, *int64, error)
