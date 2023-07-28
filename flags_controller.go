package main

import (
	"encoding/json"
	"io"
	"net/http"

	"caroline-weisberg.fun/iljournalierserver/errors"
	"caroline-weisberg.fun/iljournalierserver/services"
	"caroline-weisberg.fun/iljournalierserver/transaction"
)

type flagsController struct {
	flagsService *services.FlagsService
}

func newFlagsController(flagsService *services.FlagsService) flagsController {
	return flagsController{
		flagsService: flagsService,
	}
}

type MarkFlagRequest struct {
	UnixSeconds int64 `json:"unixSeconds"`
	FlagId      int64 `json:"flagId"`
}

type markFlagsRequestBody struct {
	accessTokenHavingObject
	Requests []MarkFlagRequest `json:"requests"`
}

func (flagsController *flagsController) markFlags(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	var markFlagsRequestBody markFlagsRequestBody
	err = json.Unmarshal(body, &markFlagsRequestBody)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if len(markFlagsRequestBody.Requests) == 0 {
		w.WriteHeader(500)
		return
	}

	markFlagsRequests := make([]transaction.MarkFlagRequest, len(markFlagsRequestBody.Requests))
	for _, request := range markFlagsRequestBody.Requests {
		markFlagsRequests = append(markFlagsRequests, transaction.MarkFlagRequest{
			UnixSeconds: request.UnixSeconds,
			FlagId:      request.FlagId,
		})
	}

	err = flagsController.flagsService.MarkFlags(r.Context(), markFlagsRequestBody.AccessToken, markFlagsRequests)
	if err != nil {
		if errors.Is(err, errors.UserNotFoundForAccessToken) {
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

	w.WriteHeader(200)
}

type addKnownFlagsRequestBody struct {
	accessTokenHavingObject
	NewFlags []string `json:"newFlags"`
}

func (flagsController *flagsController) addKnownFlags(w http.ResponseWriter, r *http.Request) {
	var body, err = io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	var addKnownFlagsRequestBody addKnownFlagsRequestBody
	err = json.Unmarshal(body, &addKnownFlagsRequestBody)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if len(addKnownFlagsRequestBody.NewFlags) == 0 {
		w.WriteHeader(500)
		return
	}

	err = flagsController.flagsService.AddKnownFlags(r.Context(), addKnownFlagsRequestBody.AccessToken, addKnownFlagsRequestBody.NewFlags)
	if err != nil {
		handleServiceError(err, w, r)
		return
	}

	w.WriteHeader(200)
}
