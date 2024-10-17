package utilities

import (
	"fmt"
	"io"
	"net/http"

	"github.com/matttm/spoticli/spoticli-cli/internal/config"
)

func GetBytesBackend(args ...interface{}) ([]byte, error) {
	var t interface{} = config.SERVER_URL
	slice := make([]interface{}, 1)
	slice[0] = t
	args = append(slice, args)
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

func getClient() *http.Client {
	return new(http.Client)
}
