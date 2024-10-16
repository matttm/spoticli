package handler

import "net/http"

func DownloadSong(title string) error {
	req, err := http.NewRequest(
		http.MethodGet,
		""
	)
	if err != nil {
		panic(err)
	}
	req.Body
}
