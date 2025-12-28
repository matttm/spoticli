package services

import "github.com/coder/flog"

type MediaService struct {
}

var mediaService *MediaService

func GetMusiceService() *MediaService {
	if mediaService == nil {
		mediaService = &MediaService{}
		println("MusicService Instantiated")
	}
	return mediaService
}

// PartitionMp3Frames takes an entire
// mp3 file and returns a slice of frames
func PartitionMp3Frames(b []byte) [][]byte {
	if len(b) == 0 {
		return [][]byte{}
	}
	flog.Infof("Total number of bytes %d", len(b))
	var frames [][]byte
	for {
		currentFrameLength := getCurrentFrameLength(b)
		if currentFrameLength <= 0 {
			break
		}
		frameData := b[:currentFrameLength]
		frames = append(frames, frameData)
		b = b[currentFrameLength:]
		flog.Infof("Frames counted %d remaining byte count %d ", len(frames), len(b))
	}
	return frames
}
