package main

import (
	"context"

	"caroline-weisberg.fun/iljournalierserver/errors"
	"caroline-weisberg.fun/iljournalierserver/services"
	"caroline-weisberg.fun/iljournalierserver/transaction"
	gojason "github.com/isaac-weisberg/go-jason"
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
	gojason.Decodable

	unixSeconds int64
	flagId      int64
}

type markFlagsRequestBody struct {
	gojason.Decodable

	accessTokenHavingRequest
	requests []MarkFlagRequest
}

func (flagsController *flagsController) markFlags(ctx context.Context, markFlagsRequestBody *markFlagsRequestBody) error {
	if len(markFlagsRequestBody.requests) == 0 {
		return errors.E("mark flags request body had no mark requests")
	}

	markFlagsRequests := make([]transaction.MarkFlagRequest, 0, len(markFlagsRequestBody.requests))
	for _, request := range markFlagsRequestBody.requests {
		markFlagsRequests = append(markFlagsRequests, transaction.MarkFlagRequest{
			UnixSeconds: request.unixSeconds,
			FlagId:      request.flagId,
		})
	}

	err := flagsController.flagsService.MarkFlags(ctx, markFlagsRequestBody.accessToken, markFlagsRequests)
	if err != nil {
		return errors.J(err, "mark flags failed")
	}

	return nil
}

type addKnownFlagsRequestBody struct {
	gojason.Decodable

	accessTokenHavingRequest
	newFlags []string
}

type addKnownFlagsResponseBody struct {
	FlagIds []int64 `json:"flagIds" validate:"required"`
}

func (flagsController *flagsController) addKnownFlags(ctx context.Context, addKnownFlagsRequestBody *addKnownFlagsRequestBody) (*addKnownFlagsResponseBody, error) {
	if len(addKnownFlagsRequestBody.newFlags) == 0 {
		return nil, errors.E("no new flags are suggested")
	}

	flagIds, err := flagsController.flagsService.AddKnownFlags(
		ctx,
		addKnownFlagsRequestBody.accessToken,
		addKnownFlagsRequestBody.newFlags,
	)
	if err != nil {
		return nil, errors.J(err, "add known flags service failed")
	}

	var addKnownFlagsResponseBody = addKnownFlagsResponseBody{
		FlagIds: *flagIds,
	}

	return &addKnownFlagsResponseBody, err
}

type getKnownFlagsRequestBody struct {
	gojason.Decodable

	accessTokenHavingRequest
}

type getKnownFlagsResponseBodyFlag struct {
	Id   *int64 `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

func newGetKnownFlagsResponseBodyFlag(id int64, name string) getKnownFlagsResponseBodyFlag {
	return getKnownFlagsResponseBodyFlag{
		Id:   &id,
		Name: name,
	}
}

type getKnownFlagsResponseBody struct {
	Flags []getKnownFlagsResponseBodyFlag `json:"flags" validate:"required"`
}

func (flagsController *flagsController) getKnownFlags(ctx context.Context, request *getKnownFlagsRequestBody) (*getKnownFlagsResponseBody, error) {
	flagModels, err := flagsController.flagsService.GetKnownFlagsForUser(ctx, request.accessToken)
	if err != nil {
		return nil, errors.J(err, "controller getKnownFlags for user failed")
	}

	var flagResponseModels = make([]getKnownFlagsResponseBodyFlag, 0, len(*flagModels))

	for _, flagModel := range *flagModels {
		flagResponseModels = append(flagResponseModels, newGetKnownFlagsResponseBodyFlag(flagModel.Id, flagModel.Name))
	}

	return &getKnownFlagsResponseBody{
		Flags: flagResponseModels,
	}, nil
}
