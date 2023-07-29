package main

import (
	"context"

	"caroline-weisberg.fun/iljournalierserver/errors"
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
	UnixSeconds int64  `json:"unixSeconds"`
	Msg         string `json:"msg"`
}

func (moreMessagesController *moreMessagesController) addMoreMessage(ctx context.Context, addMoreMessageRequestBody *addMoreMessageRequestBody) error {
	var err = moreMessagesController.moreMessagesService.AddMessage(
		ctx,
		addMoreMessageRequestBody.AccessToken,
		addMoreMessageRequestBody.UnixSeconds,
		addMoreMessageRequestBody.Msg,
	)

	if err != nil {
		return errors.J(err, "add message failed")
	}

	return nil
}
