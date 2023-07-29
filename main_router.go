package main

import (
	"encoding/json"
	"io"
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

	var statusCode int
	var body *[]byte

	switch r.Method {
	case http.MethodPost:
		switch inAppRoute {
		case "/user/create":
			router.userController.createUser(w, r)
		case "/messages/add":
			router.moreMessagesController.addMoreMessage(w, r)
		case "/flags/addflags":
			statusCode, body = router.addKnownFlags(r)
		case "/flags/mark":
			statusCode = router.markFlags(r)
		default:
			statusCode = 404
		}
	default:
		statusCode = 404
	}

	w.WriteHeader(statusCode)
	if body != nil {
		w.Write(*body)
	}
}

func (router *mainRouter) markFlags(r *http.Request) int {
	err := router.flagsController.markFlags(r)
	if err != nil {
		return router.handleCommonErrors(err)
	}

	return 200
}

func (router *mainRouter) addKnownFlags(r *http.Request) (int, *[]byte) {
	var body, err = io.ReadAll(r.Body)
	if err != nil {
		return 500, nil
	}

	var addKnownFlagsRequestBody addKnownFlagsRequestBody
	err = json.Unmarshal(body, &addKnownFlagsRequestBody)
	if err != nil {
		return 500, nil
	}

	addKnownFlagsResponseBody, err := router.flagsController.addKnownFlags(r.Context(), addKnownFlagsRequestBody)
	if err != nil {
		return router.handleCommonErrors(err), nil
	}

	responseBody, err := json.Marshal(addKnownFlagsResponseBody)
	if err != nil {
		return 500, nil
	}

	return 200, &responseBody
}

func (router *mainRouter) respond404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
}

func handleRoute(route func(w http.ResponseWriter, r http.Request) (int, *[]byte)) {

}

func (router *mainRouter) handleCommonErrors(err error) int {
	if errors.Is(err, errors.UserNotFoundForAccessToken) {
		return 418
	} else if errors.Is(err, errors.UserNotFoundForMagicKey) {
		return 418
	} else if errors.Is(err, errors.FlagDoesntBelongToTheUser) {
		return 418
	} else {
		return 500
	}
}
