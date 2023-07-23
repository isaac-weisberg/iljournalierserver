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
	inAppRoute, found := strings.CutPrefix(r.URL.Path, "/iljournalierserver")

	if !found {
		router.respond404(w, r)
		return
	}

	switch r.Method {
	case http.MethodPost:
		switch inAppRoute {
		case "/user/create":
			router.userController.createUser(w, r)
		default:
			router.respond404(w, r)
		}
	default:
		router.respond404(w, r)
	}
}

func (router *mainRouter) respond404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
}
