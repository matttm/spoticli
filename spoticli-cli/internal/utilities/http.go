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
	args = append(slice, args...)
	fmt.Println(fmt.Sprintf("http://%s/%s", args...))
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("http://%s/%s", args...),
		nil,
	)
	if err != nil {
		panic(err)
	}
	res, err := getClient().Do(req)
	if err != nil {
		panic(err)
	}
	var data []byte
	data, err = io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	return data, nil
}

func getClient() *http.Client {
	return new(http.Client)
}
