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
	t.FileSize = 5000000
	return t, nil
}
