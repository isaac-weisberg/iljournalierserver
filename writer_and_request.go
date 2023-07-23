package main

import "net/http"

type writerAndRequest struct {
	w http.ResponseWriter
	r *http.Request
}
