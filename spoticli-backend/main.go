package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mstttm/spoticli/spoticli-backend/internal/models"
	"github.com/mstttm/spoticli/spoticli-backend/internal/routers"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", HealthHandler)
	r.Handle("/audio", routers.AudioRouter())

	http.Handle("/", r)
	fmt.Println("Listening on locslhost:4200")
	printEndpoints(r)
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

func printEndpoints(r *mux.Router) {
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}
		methods, err := route.GetMethods()
		if err != nil {
			return nil
		}
		fmt.Printf("%v %s\n", methods, path)
		return nil
	})
}
