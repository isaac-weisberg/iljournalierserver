package main

import (
	"context"

	"caroline-weisberg.fun/iljournalierserver/errors"
	"caroline-weisberg.fun/iljournalierserver/models"
	"caroline-weisberg.fun/iljournalierserver/services"
	gojason "github.com/isaac-weisberg/go-jason"
)

type moreMessagesController struct {
	moreMessagesService *services.MoreMessagesService
}

func newMoreMessagesController(moreMessagesService *services.MoreMessagesService) moreMessagesController {
	return moreMessagesController{moreMessagesService: moreMessagesService}
}

type addMoreMessageRequest struct {
	gojason.Decodable

	unixSeconds int64
	msg         string
}

type addMoreMessageRequestBody struct {
	gojason.Decodable

	accessTokenHavingRequest
	requests []addMoreMessageRequest
}

func (moreMessagesController *moreMessagesController) addMoreMessages(
	ctx context.Context,
	addMoreMessageRequestBody *addMoreMessageRequestBody,
) error {
	if len(addMoreMessageRequestBody.requests) == 0 {
		return errors.E("can't insert more messages without more message requests")
	}

	var addMoreMessageRequests = make([]models.AddMessageRequest, 0, len(addMoreMessageRequestBody.requests))
	for _, request := range addMoreMessageRequestBody.requests {
		addMoreMessageRequests = append(addMoreMessageRequests, models.NewAddMessageRequest(*&request.unixSeconds, request.msg))
	}

	var err = moreMessagesController.moreMessagesService.AddMessage(
		ctx,
		addMoreMessageRequestBody.accessToken,
		addMoreMessageRequests,
	)

	if err != nil {
		return errors.J(err, "add message failed")
	}

	return nil
}
