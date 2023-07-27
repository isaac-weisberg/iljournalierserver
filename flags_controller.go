package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type flagsController struct {
	flagsService *flagsService
}

func newFlagsController(flagsService *flagsService) flagsController {
	return flagsController{
		flagsService: flagsService,
	}
}

type markFlagRequest struct {
	UnixSeconds int64 `json:"unixSeconds"`
	FlagId      int64 `json:"flagId"`
}

type markFlagsRequestBody struct {
	accessTokenHavingObject
	Requests []markFlagRequest `json:"requests"`
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

	err = flagsController.flagsService.markFlags(r.Context(), markFlagsRequestBody.AccessToken, markFlagsRequestBody.Requests)
	if err != nil {
		if is(err, userNotFoundForAccessToken) {
			w.WriteHeader(418)
			return
		} else if is(err, flagDoesntBelongToTheUser) {
			w.WriteHeader(418)
			return
		} else {
			w.WriteHeader(500)
			return
		}
	}

	w.WriteHeader(200)
}
