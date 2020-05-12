package main

import "net/http"

// renders home view with a search line
func homeHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {}
}

// runs all shorts, if not found skip
func runURLsHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {}
}

// render settings view with a search line and table of all urls
func homeSettingsHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {}
}

// create a new url
func cURLHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {}
}

// update an url
func uURLHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {}
}

// delete an url
func dURLHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {}
}
