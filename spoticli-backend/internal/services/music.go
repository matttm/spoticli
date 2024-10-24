package services

import (
	"github.com/matttm/spoticli/spoticli-models"
)

type MusicService struct {
}

var musicService *MusicService

func GetMusiceService() *MusicService {
	if musicService == nil {
		musicService = &MusicService{}
		println("MusicService Instantiated")
	}
	return musicService
}
func GetTrack(id int) (*models.Track, error) {
	t := new(models.Track)
	t.Title = "bat_country.mp3"
	t.FileSize = 7523726
	return t, nil
}
