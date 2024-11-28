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
// mp3 file andreturns a slice of frames
func PartitionMp3Frames(b []byte) [][]byte {
	if len(b) == 0 {
		return [][]byte{}
	}
	var frames [][]byte
	for {
		currentFrameLength := getCurrentFrameLength(b)
		clip := b[:currentFrameLength]
		frames = append(frames, clip)
		b = b[currentFrameLength:]
		if currentFrameLength <= 0 {
			break
		}
		flog.Infof("Frame count %d frame length %d ", len(frames), currentFrameLength)
	}
	return frames
}
