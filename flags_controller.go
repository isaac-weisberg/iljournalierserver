package main

import (
	"context"
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

func (flagsController *flagsController) addKnownFlags(ctx context.Context, addKnownFlagsRequestBody addKnownFlagsRequestBody) (*addKnownFlagsResponseBody, error) {
	if len(addKnownFlagsRequestBody.NewFlags) == 0 {
		return nil, errors.E("no new flags are suggested")
	}

	flagIds, err := flagsController.flagsService.AddKnownFlags(
		ctx,
		addKnownFlagsRequestBody.AccessToken,
		addKnownFlagsRequestBody.NewFlags,
	)
	if err != nil {
		return nil, errors.J(err, "add known flags service failed")
	}

	var addKnownFlagsResponseBody = addKnownFlagsResponseBody{
		FlagIds: *flagIds,
	}

	return &addKnownFlagsResponseBody, err
}
