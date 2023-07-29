package main

import (
	"context"
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
			statusCode, body = router.markFlags(r)
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

func handleWJsonReq[
	ReqBody any,
](
	router *mainRouter,
	r *http.Request,
	handler func(ctx context.Context, body *ReqBody) (int, *[]byte),
) (int, *[]byte) {
	var body, err = io.ReadAll(r.Body)
	if err != nil {
		return 500, nil
	}

	var parsedBody ReqBody
	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		return 500, nil
	}

	return handler(r.Context(), &parsedBody)
}

func handleWJsonReqJsonRes[
	ReqBody any,
	ResBody any,
](
	router *mainRouter,
	r *http.Request,
	handler func(ctx context.Context, body *ReqBody) (int, *ResBody),
) (int, *[]byte) {
	return handleWJsonReq[ReqBody](router, r, func(ctx context.Context, reqBody *ReqBody) (int, *[]byte) {
		statusCode, resBody := handler(ctx, reqBody)

		if resBody == nil {
			return statusCode, nil
		}

		bodyBytes, err := json.Marshal(resBody)
		if err != nil {
			return 500, nil
		}

		return statusCode, &bodyBytes
	})
}

func (router *mainRouter) markFlags(r *http.Request) (int, *[]byte) {
	return handleWJsonReq[markFlagsRequestBody](router, r, func(ctx context.Context, body *markFlagsRequestBody) (int, *[]byte) {
		err := router.flagsController.markFlags(ctx, body)
		if err != nil {
			return router.handleCommonErrors(err), nil
		}

		return 200, nil
	})
}

func (router *mainRouter) addKnownFlags(r *http.Request) (int, *[]byte) {
	return handleWJsonReqJsonRes[addKnownFlagsRequestBody, addKnownFlagsResponseBody](router, r, func(ctx context.Context, body *addKnownFlagsRequestBody) (int, *addKnownFlagsResponseBody) {
		addKnownFlagsResponseBody, err := router.flagsController.addKnownFlags(r.Context(), body)
		if err != nil {
			return router.handleCommonErrors(err), nil
		}

		return 200, addKnownFlagsResponseBody
	})
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
