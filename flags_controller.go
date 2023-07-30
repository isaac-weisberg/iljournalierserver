package main

import (
	"context"

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
	UnixSeconds *int64 `json:"unixSeconds" validate:"required"`
	FlagId      *int64 `json:"flagId" validate:"required"`
}

type markFlagsRequestBody struct {
	accessTokenHavingObject
	Requests []MarkFlagRequest `json:"requests" validate:"required"`
}

func (flagsController *flagsController) markFlags(ctx context.Context, markFlagsRequestBody *markFlagsRequestBody) error {
	if len(markFlagsRequestBody.Requests) == 0 {
		return errors.E("mark flags request body had no mark requests")
	}

	markFlagsRequests := make([]transaction.MarkFlagRequest, 0, len(markFlagsRequestBody.Requests))
	for _, request := range markFlagsRequestBody.Requests {
		markFlagsRequests = append(markFlagsRequests, transaction.MarkFlagRequest{
			UnixSeconds: *request.UnixSeconds,
			FlagId:      *request.FlagId,
		})
	}

	err := flagsController.flagsService.MarkFlags(ctx, markFlagsRequestBody.AccessToken, markFlagsRequests)
	if err != nil {
		return errors.J(err, "mark flags failed")
	}

	return nil
}

type addKnownFlagsRequestBody struct {
	accessTokenHavingObject
	NewFlags []string `json:"newFlags" validate:"required"`
}

type addKnownFlagsResponseBody struct {
	FlagIds []int64 `json:"flagIds" validate:"required"`
}

func (flagsController *flagsController) addKnownFlags(ctx context.Context, addKnownFlagsRequestBody *addKnownFlagsRequestBody) (*addKnownFlagsResponseBody, error) {
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

type getKnownFlagsRequestBody struct {
	accessTokenHavingObject
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
	flagModels, err := flagsController.flagsService.GetKnownFlagsForUser(ctx, request.AccessToken)
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
