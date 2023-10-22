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
	utils.DebugLog("IlJournalierServer: Got connection!", r.URL)
	inAppRoute, found := strings.CutPrefix(r.URL.Path, "/iljournalierserver")

	if !found {
		router.respond404(w, r)
		return
	}

	var statusCode int
	var resBody *[]byte

	reqBody, err := io.ReadAll(r.Body)
	if err == nil {
		switch r.Method {
		case http.MethodPost:
			switch inAppRoute {
			case "/user/login":
				statusCode, resBody = router.handleAndConvert(router.login(r, reqBody))
			case "/user/create":
				statusCode, resBody = router.handleAndConvert(router.createUser(r))
			case "/messages/add":
				statusCode, resBody = router.handleAndConvert(router.addMoreMessage(r, reqBody))
			case "/flags/addflags":
				statusCode, resBody = router.handleAndConvert(router.addKnownFlags(r, reqBody))
			case "/flags/mark":
				statusCode, resBody = router.handleAndConvert(router.markFlags(r, reqBody))
			case "/flags":
				statusCode, resBody = router.handleAndConvert(router.getKnownFlags(r, reqBody))
			default:
				statusCode = 404
			}
		default:
			statusCode = 404
		}
	} else {
		statusCode = 500
	}

	utils.WriteDebugCorsHeader(w, r)
	w.WriteHeader(statusCode)

	if resBody != nil {
		w.Write(*resBody)
	}
}

func (router *mainRouter) login(r *http.Request, body []byte) (*[]byte, error) {
	loginRequestBody, err := makeLoginRequestBodyFromJson(body)
	if err != nil {
		return nil, errors.J(err, "parse body failed")
	}

	response, err := router.userController.login(r.Context(), loginRequestBody)
	if err != nil {
		return nil, errors.J(err, "login failed")
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		return nil, errors.J(err, "marshal response failed")
	}

	return &bytes, nil
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

func (router *mainRouter) addMoreMessage(r *http.Request, body []byte) (*[]byte, error) {
	addMoreMessageRequestBody, err := makeAddMoreMessageRequestBodyFromJson(body)
	if err != nil {
		return nil, errors.J(err, "parse json failed")
	}

	err = router.moreMessagesController.addMoreMessages(r.Context(), addMoreMessageRequestBody)
	if err != nil {
		return nil, errors.J(err, "add more message failed")
	}

	return nil, nil
}

func (router *mainRouter) markFlags(r *http.Request, body []byte) (*[]byte, error) {
	markFlagsRequestBody, err := makeMarkFlagsRequestBodyFromJson(body)
	if err != nil {
		return nil, errors.J(err, "parsing body failed")
	}

	err = router.flagsController.markFlags(r.Context(), markFlagsRequestBody)
	if err != nil {
		return nil, errors.J(err, "mark flags failed")
	}

	return nil, nil
}

func (router *mainRouter) addKnownFlags(r *http.Request, body []byte) (*[]byte, error) {
	addKnownFlagsRequestBody, err := makeAddKnownFlagsRequestBodyFromJson(body)
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

func (router *mainRouter) getKnownFlags(r *http.Request, body []byte) (*[]byte, error) {
	getKnownFlagsRequestBody, err := makeGetKnownFlagsRequestBodyFromJson(body)
	if err != nil {
		return nil, errors.J(err, "parsing request body failed")
	}

	flagModels, err := router.flagsController.getKnownFlags(r.Context(), getKnownFlagsRequestBody)
	if err != nil {
		return nil, errors.J(err, "get known flags failed")
	}

	bytes, err := json.Marshal(flagModels)
	if err != nil {
		return nil, errors.J(err, "serializing body failed")
	}

	return &bytes, nil
}

func (router *mainRouter) handleAndConvert(responseBody *[]byte, err error) (int, *[]byte) {
	if err != nil {
		utils.DebugLog(err)
		return router.handleCommonErrors(err), nil
	}
	return 200, responseBody
}

func (router *mainRouter) respond404(w http.ResponseWriter, r *http.Request) {
	utils.WriteDebugCorsHeader(w, r)
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
