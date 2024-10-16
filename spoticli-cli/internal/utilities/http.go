package utilities

import (
	"fmt"
	"io"
	"net/http"

	"github.com/matttm/spoticli/spoticli-cli/internal/config"
)

func GetBytesBackend(args ...interface{}) ([]byte, error) {
	args = append([]interface{ config.SERVER_URL }, args)
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("http://%s/%s", args...),
		nil,
	)
	if err != nil {
		panic(err)
	}
	res, err := getClient().Do(req)
	body := res.Body
	defer res.Body.Close()
	b, err := io.ReadAll(body)
	return b, nil
}

func getClient() http.Client {
	return http.Client{}
}
