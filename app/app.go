package main

import "github.com/gorilla/mux"

// run an app
func app() {
	client := initDB(mongoURL)

	// routing
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler()).Methods("GET")
	r.HandleFunc("/q/{shorts}", runURLsHandler()).Methods("POST")

	set := r.PathPrefix("/settings").Subrouter()
	set.HandleFunc("/", homeSettingsHandler()).Methods("GET")
	set.HandleFunc("/c", cURLHandler()).Methods("POST")
	set.HandleFunc("/u/{id}", uURLHandler()).Methods("POST")
	set.HandleFunc("/d/{id}", dURLHandler()).Methods("POST")

	defer client.Disconnect(dbCtx)
}
