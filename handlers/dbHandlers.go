package handlers

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"speedit/models"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/browser"
	"go.mongodb.org/mongo-driver/mongo"
)

// HomeHandler renders home view with a search line
func HomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		html, err := readHTML("home.gohtml")
		if err != nil {
			log.Fatal(err)
		}

		if _, err = w.Write(html); err != nil {
			log.Fatal(err)
		}
	}
}

// RunURLsHandler runs all shorts, if not found skip
func RunURLsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		strURLs := parseURLs(strings.TrimSpace(r.FormValue("shorts")))

		if strURLs == nil {
			http.Redirect(w, r, "https://www.google.com", http.StatusSeeOther)
			return
		}

		// runs each url if found
		for _, strURL := range strURLs {

			// else skip
			u, err := ReadOneURL(strURL)
			if err == mongo.ErrNoDocuments {
				continue
			} else if err != nil {
				log.Fatal(err)
			}

			if err := browser.OpenURL(u.Origin); err != nil {
				log.Fatal(err)
			}
		}

		html, err := readHTML("close.gohtml")
		if err != nil {
			log.Fatal(err)
		}

		var t = template.Must(template.New("closing").Parse(string(html)))
		t.Execute(w, nil)

		/*
			if _, err = w.Write(html); err != nil {
				/log.Fatal(err)
			}
		*/
	}
}

// HomeSettingsHandler render settings view with a search line and table of all urls
func HomeSettingsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urls, err := ReadAllURLs()
		if err != nil {
			log.Fatal(err)
		}

		// TODO plush and set it
		for _, u := range urls {
			fmt.Fprintln(w, u)
		}
	}
}

// CURLHandler create a new url
func CURLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := models.Url{
			ID:     uuid.New().String(),
			Origin: r.FormValue("long"),
			Short:  r.FormValue("short"),
			Usage:  0,
		}

		if err := CreateURL(u); err != nil {
			log.Fatal(err)
		}
	}
}

// UURLHandler update an url
func UURLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		u := models.Url{
			ID:     id,
			Origin: r.FormValue("long"),
			Short:  r.FormValue("short"),
			Usage:  0,
		}

		if err := UpdateURL(id, u); err != nil {
			log.Fatal(err)
		}
	}
}

// DURLHandler delete an url
func DURLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		DeleteURL(mux.Vars(r)["id"])
	}
}

// split string
func parseURLs(strURLs string) []string {
	if strURLs == "" {
		return nil
	}

	return strings.Split(strURLs, " ")
}

// get template
func readHTML(gohtml string) ([]byte, error) {
	path := "templates/" + gohtml
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}, fmt.Errorf("Could not read from file %v, error: %v", path, err)
	}

	return data, nil
}
