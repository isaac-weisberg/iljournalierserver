package main

import (
	"net/http"
)

func main() {
	diContainer, err := NewDIContainer()

	if err != nil {
		panic(err)
	}

	mainRouter := NewMainRouter(diContainer)

	mux := http.NewServeMux()

	mux.Handle("/", &mainRouter)

	http.ListenAndServe(":8081", mux)
}
