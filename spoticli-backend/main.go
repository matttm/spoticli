package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	datbase "github.com/matttm/spoticli/spoticli-backend/internal/database"
	"github.com/matttm/spoticli/spoticli-backend/internal/routers"
	"github.com/matttm/spoticli/spoticli-models"
	"github.com/rs/cors"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", HealthHandler)
	routers.AttachAudioRouter(r.PathPrefix("/audio").Subrouter())
	routers.AttachFileMetaInfoRouter(r.PathPrefix("/files").Subrouter())

	datbase.InitializeDatabase()

	http.Handle("/", r)
	fmt.Println("Listening on localhost:4200")

	// adding cors for null origin testing
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)
	printEndpoints(r)
	err := http.ListenAndServe(":4200", handler)
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
