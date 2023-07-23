package main

import (
	"net/http"
	"strings"
)

type MainRouter struct {
	userController *UserController
}

func (router *MainRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writerAndRequest := WriterAndRequest{w, r}

	inAppRoute, found := strings.CutPrefix(r.URL.Path, "/iljournalierserver")

	if !found {
		router.Respond404(writerAndRequest)
		return
	}

	switch r.Method {
	case http.MethodPost:
		switch inAppRoute {
		case "/user/create":
			router.userController.createUser(writerAndRequest)
		default:
			router.Respond404(writerAndRequest)
		}
	default:
		router.Respond404(writerAndRequest)
	}
}

func (router *MainRouter) Respond404(writerAndRequest WriterAndRequest) {
	writerAndRequest.w.WriteHeader(404)
}
