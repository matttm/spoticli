package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mstttm/spoticli/spoticli-backend/internal/models"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HealthHandler)
	http.Handle("/", r)
	fmt.Println("Listening on locslhost:4200")
	err := http.ListenAndServe(":4200", r)
	if err != nil {
		panic(err)
	}
}
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	hc := models.HealthCheck{
		ServiceName: "spoticli-backend",
		Version:     "v1.0.0",
		ManagedBy:   "matttm",
		Timestamp:   time.Now().String(),
	}
	b, err := json.Marshal(hc)
	if err != nil {
		panic(err)
	}
	w.Write(b)
}
