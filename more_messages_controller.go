package main

import (
	"context"

	"caroline-weisberg.fun/iljournalierserver/errors"
	"caroline-weisberg.fun/iljournalierserver/models"
	"caroline-weisberg.fun/iljournalierserver/services"
)

type moreMessagesController struct {
	moreMessagesService *services.MoreMessagesService
}

func newMoreMessagesController(moreMessagesService *services.MoreMessagesService) moreMessagesController {
	return moreMessagesController{moreMessagesService: moreMessagesService}
}

type addMoreMessageRequestBody struct {
	accessTokenHavingObject
	Requests []struct {
		UnixSeconds *int64 `json:"unixSeconds" validate:"required"`
		Msg         string `json:"msg" validate:"required"`
	} `json:"requests" validate:"required"`
}

func (moreMessagesController *moreMessagesController) addMoreMessages(
	ctx context.Context,
	addMoreMessageRequestBody *addMoreMessageRequestBody,
) error {
	if len(addMoreMessageRequestBody.Requests) == 0 {
		return errors.E("can't insert more messages without more message requests")
	}

	var addMoreMessageRequests = make([]models.AddMessageRequest, 0, len(addMoreMessageRequestBody.Requests))
	for _, request := range addMoreMessageRequestBody.Requests {
		addMoreMessageRequests = append(addMoreMessageRequests, models.NewAddMessageRequest(*request.UnixSeconds, request.Msg))
	}

	var err = moreMessagesController.moreMessagesService.AddMessage(
		ctx,
		addMoreMessageRequestBody.AccessToken,
		addMoreMessageRequests,
	)

	if err != nil {
		return errors.J(err, "add message failed")
	}

	return nil
}
