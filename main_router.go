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
			statusCode, body = router.handleAndConvert(router.createUser(r))
		case "/messages/add":
			statusCode, body = router.handleAndConvert(router.addMoreMessage(r))
		case "/flags/addflags":
			statusCode, body = router.handleAndConvert(router.addKnownFlags(r))
		case "/flags/mark":
			statusCode, body = router.handleAndConvert(router.markFlags(r))
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

func parseJson[R any](input io.ReadCloser) (*R, error) {
	var body, err = io.ReadAll(input)
	if err != nil {
		return nil, errors.J(err, "failed parsing")
	}

	var parsedBody R
	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		return nil, errors.J(err, "parsing json failed")
	}

	return &parsedBody, nil
}

func (router *mainRouter) createUser(r *http.Request) (*[]byte, error) {
	response, err := router.userController.createUser(r.Context())
	if err != nil {
		return nil, errors.J(err, "create user controller failed")
	}

	body, err := json.Marshal(response)
	if err != nil {
		return nil, errors.J(err, "body serialization failed")
	}

	return &body, nil
}

func (router *mainRouter) addMoreMessage(r *http.Request) (*[]byte, error) {
	addMoreMessageRequestBody, err := parseJson[addMoreMessageRequestBody](r.Body)
	if err != nil {
		return nil, errors.J(err, "parse json failed")
	}

	err = router.moreMessagesController.addMoreMessage(r.Context(), addMoreMessageRequestBody)
	if err != nil {
		return nil, errors.J(err, "add more message failed")
	}

	return nil, nil
}

func (router *mainRouter) markFlags(r *http.Request) (*[]byte, error) {
	markFlagsRequestBody, err := parseJson[markFlagsRequestBody](r.Body)
	if err != nil {
		return nil, errors.J(err, "parsing body failed")
	}

	err = router.flagsController.markFlags(r.Context(), markFlagsRequestBody)
	if err != nil {
		return nil, errors.J(err, "mark flags failed")
	}

	return nil, nil
}

func (router *mainRouter) addKnownFlags(r *http.Request) (*[]byte, error) {
	addKnownFlagsRequestBody, err := parseJson[addKnownFlagsRequestBody](r.Body)
	if err != nil {
		return nil, errors.J(err, "parsing body failed")
	}

	addKnownFlagsResponseBody, err := router.flagsController.addKnownFlags(r.Context(), addKnownFlagsRequestBody)
	if err != nil {
		return nil, errors.J(err, "add known flags failed")
	}

	responseBytes, err := json.Marshal(addKnownFlagsResponseBody)
	if err != nil {
		return nil, errors.J(err, "serializing body failed")
	}

	return &responseBytes, nil
}

func (router *mainRouter) handleAndConvert(responseBody *[]byte, err error) (int, *[]byte) {
	if err != nil {
		return router.handleCommonErrors(err), nil
	}
	return 200, responseBody
}

func (router *mainRouter) respond404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
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
