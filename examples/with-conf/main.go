package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/noelukwa/tempest"
)

//go:embed views
var views embed.FS

func main() {

	templates, err := tempest.WithConfig(&tempest.Config{
		IncludesDir: "partials",
		Layout:      "base",
	}).LoadFS(views)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		home := templates["home"]
		if home == nil {
			http.Error(w, "template not found", http.StatusInternalServerError)
			return
		}
		home.Execute(w, nil)
	})

	log.Fatal(http.ListenAndServe(":8087", mux))
}
