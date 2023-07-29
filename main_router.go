package main

import (
	"net/http"
	"strings"

	"caroline-weisberg.fun/iljournalierserver/utils"
)

type mainRouter struct {
	userController         *userController
	moreMessagesController *moreMessagesController
	flagsController        *flagsController
}

func newMainRouter(di *diContainer) mainRouter {
	userController := newUserController(di.userService)
	moreMessagesController := newMoreMessagesController(di.moreMessagesService)
	flagsController := newFlagsController(di.flagsService)

	mainRouter := mainRouter{
		userController:         &userController,
		moreMessagesController: &moreMessagesController,
		flagsController:        &flagsController,
	}

	return mainRouter
}

func (router *mainRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	utils.Log("IlJournalierServer: Got connection!", r.URL)
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
		case "/flags/addflags":
			router.flagsController.addKnownFlags(w, r)
		case "/flags/mark":
			router.flagsController.markFlags(w, r)
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
