package main

import (
	"context"
	"net/http"
	"speedit/handlers"
	"speedit/models"

	"github.com/gorilla/mux"
)

func main() {
	ctx := context.Background()
	client := handlers.InitDB(ctx, models.MongoURL)
	defer client.Disconnect(ctx)

	// routing
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.HomeHandler()).Methods("GET")
	r.HandleFunc("/", handlers.RunURLsHandler()).Methods("POST")

	set := r.PathPrefix("/set").Subrouter()
	set.HandleFunc("/", handlers.HomeSettingsHandler()).Methods("GET")
	set.HandleFunc("/c", handlers.CURLHandler()).Methods("POST")
	set.HandleFunc("/u/{id}", handlers.UURLHandler()).Methods("POST")
	set.HandleFunc("/d/{id}", handlers.DURLHandler()).Methods("POST")

	http.ListenAndServe(":8080", r)
}
