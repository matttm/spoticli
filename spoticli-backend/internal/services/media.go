package services

import (
	"bytes"
	"encoding/hex"
	"io"

	"github.com/matttm/spoticli/spoticli-models"
)

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
func GetTrack(id int) (*models.Track, error) {
	t := new(models.Track)
	t.Title = "bat_country.mp3"
	t.FileSize = 7523726
	return t, nil
}

func ReadID3vsHeader(r io.Reader) []byte {
	identifier := make([]byte, 3)
	n, err := r.Read(identifier)
	if err != nil {
	}
	idString := "ID3"
	decodedId, err := hex.DecodeString(idString)
	if err != nil {
	}
	if !bytes.Equal(identifier, decodedId) {
		panic("No ID3 indicator found")
	}
}
