package main

import (
	"net/http"

	"crawler/crawler_single/frontend/controller"
)

func main() {
	http.Handle("/", http.FileServer(
		http.Dir("crawler_single/frontend/view")))
	http.Handle("/search",
		controller.CreateSearchResultHandler(
			"crawler_single/frontend/view/template.html"))
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic(err)
	}
}
