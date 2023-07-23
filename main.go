package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mainRouter := MainRouter{}

	mux.Handle("/", &mainRouter)

	http.ListenAndServe(":8081", mux)
}
