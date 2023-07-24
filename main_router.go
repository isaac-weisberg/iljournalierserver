package main

import (
	"net/http"
	"strings"
)

type mainRouter struct {
	userController         *userController
	moreMessagesController *moreMessagesController
}

func newMainRouter(di *diContainer) mainRouter {
	userController := newUserController(di.userService)
	moreMessagesController := newMoreMessagesController(di.moreMessagesService)

	mainRouter := mainRouter{
		userController:         &userController,
		moreMessagesController: &moreMessagesController,
	}

	return mainRouter
}

func (router *mainRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log("IlJournalierServer: Got connection!", r.URL)
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
		case "/messages/add":
			router.moreMessagesController.addMoreMessage(w, r)
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
