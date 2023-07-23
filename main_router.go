package main

import (
	"net/http"
	"strings"
)

type mainRouter struct {
	userController *userController
}

func newMainRouter(di *diContainer) mainRouter {
	userController := newUserController(&di.userService)

	mainRouter := mainRouter{userController: &userController}

	return mainRouter
}

func (router *mainRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writerAndRequest := writerAndRequest{w, r}

	inAppRoute, found := strings.CutPrefix(r.URL.Path, "/iljournalierserver")

	if !found {
		router.respond404(writerAndRequest)
		return
	}

	switch r.Method {
	case http.MethodPost:
		switch inAppRoute {
		case "/user/create":
			router.userController.createUser(writerAndRequest)
		default:
			router.respond404(writerAndRequest)
		}
	default:
		router.respond404(writerAndRequest)
	}
}

func (router *mainRouter) respond404(writerAndRequest writerAndRequest) {
	writerAndRequest.w.WriteHeader(404)
}
