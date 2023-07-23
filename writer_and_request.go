package main

import "net/http"

type WriterAndRequest struct {
	w http.ResponseWriter
	r *http.Request
}
