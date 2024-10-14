package services

import "github.com/mstttm/spoticli/spoticli-backend/internal/models"

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
	t.Title = "bat-country.mp3"
	return t, nil
}
