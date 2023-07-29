package main

import (
	"encoding/json"
	"io"
	"net/http"

	"caroline-weisberg.fun/iljournalierserver/errors"
	"caroline-weisberg.fun/iljournalierserver/services"
	"caroline-weisberg.fun/iljournalierserver/transaction"
	"caroline-weisberg.fun/iljournalierserver/utils"
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

func (flagsController *flagsController) markFlags(r *http.Request) error {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		return errors.J(err, "read body failed")
	}

	var markFlagsRequestBody markFlagsRequestBody
	err = json.Unmarshal(body, &markFlagsRequestBody)
	if err != nil {
		return errors.J(err, "parse body failed")
	}

	if len(markFlagsRequestBody.Requests) == 0 {
		return errors.E("mark flags request body had no mark requests")
	}

	markFlagsRequests := make([]transaction.MarkFlagRequest, 0, len(markFlagsRequestBody.Requests))
	for _, request := range markFlagsRequestBody.Requests {
		markFlagsRequests = append(markFlagsRequests, transaction.MarkFlagRequest{
			UnixSeconds: request.UnixSeconds,
			FlagId:      request.FlagId,
		})
	}

	err = flagsController.flagsService.MarkFlags(r.Context(), markFlagsRequestBody.AccessToken, markFlagsRequests)
	if err != nil {
		return errors.J(err, "mark flags failed")
	}
	return nil
}

type addKnownFlagsRequestBody struct {
	accessTokenHavingObject
	NewFlags []string `json:"newFlags"`
}

type addKnownFlagsResponseBody struct {
	FlagIds []int64 `json:"flagIds"`
}

func (flagsController *flagsController) addKnownFlags(w http.ResponseWriter, r *http.Request) {
	var body, err = io.ReadAll(r.Body)
	if err != nil {
		utils.Log(err)
		w.WriteHeader(500)
		return
	}

	var addKnownFlagsRequestBody addKnownFlagsRequestBody
	err = json.Unmarshal(body, &addKnownFlagsRequestBody)
	if err != nil {
		utils.Log(err)
		w.WriteHeader(500)
		return
	}

	if len(addKnownFlagsRequestBody.NewFlags) == 0 {
		w.WriteHeader(500)
		return
	}

	flagIds, err := flagsController.flagsService.AddKnownFlags(
		r.Context(),
		addKnownFlagsRequestBody.AccessToken,
		addKnownFlagsRequestBody.NewFlags,
	)
	if err != nil {
		utils.Log(err)
		return
	}

	var addKnownFlagsResponseBody = addKnownFlagsResponseBody{
		FlagIds: *flagIds,
	}

	responseBody, err := json.Marshal(addKnownFlagsResponseBody)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	w.Write(responseBody)
}
