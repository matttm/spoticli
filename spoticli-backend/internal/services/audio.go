package services

import (
	"io"
)

func GetPresignedUrl(key string) (string, error) {
	svc := GetStorageService()
	req, err := svc.GetPresignedUrl("bat_country.mp3")
	if err != nil {
		panic(err)
	}
	return req.URL, nil
}
func GetAudio(key string) ([]byte, error) {
	svc := GetStorageService()
	res, err := svc.DownloadFile("bat_country.mp3")
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}
