package main

import (
	"net/http"
	"strings"

	"caroline-weisberg.fun/iljournalierserver/errors"
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
			router.markFlags(w, r)
		default:
			router.respond404(w, r)
		}
	default:
		router.respond404(w, r)
	}
}

func (router *mainRouter) markFlags(w http.ResponseWriter, r *http.Request) {
	err := router.flagsController.markFlags(r)
	if err != nil {
		router.handleCommonErrors(err, w)
		return
	}

	w.WriteHeader(200)
}

func (router *mainRouter) respond404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
}

func (router *mainRouter) handleCommonErrors(err error, w http.ResponseWriter) {
	if errors.Is(err, errors.UserNotFoundForAccessToken) {
		w.WriteHeader(418)
		return
	} else if errors.Is(err, errors.UserNotFoundForMagicKey) {
		w.WriteHeader(418)
		return
	} else if errors.Is(err, errors.FlagDoesntBelongToTheUser) {
		w.WriteHeader(418)
		return
	} else {
		w.WriteHeader(500)
		return
	}
}
