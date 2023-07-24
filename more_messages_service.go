package main

import (
	"context"
)

type moreMessagesService struct {
	databaseService *databaseService
}

func newMoreMessagesService(databaseService *databaseService) moreMessagesService {
	return moreMessagesService{databaseService: databaseService}
}

var userNotFoundForAccessToken = e("userNotFoundForAccessToken")

func (moreMessagesService *moreMessagesService) addMessage(ctx context.Context, accessToken string, msg string) error {
	tx, err := moreMessagesService.databaseService.beginTx(ctx)

	if err != nil {
		return j(err, "tx create failed")
	}

	userId, err := tx.findUserIdForAccessToken(accessToken)
	if err != nil {
		return j(err, "find user for accessToken failed")
	}

	if userId == nil {
		return userNotFoundForAccessToken
	}

	err = tx.addMoreMessage(*userId, msg)
	if err != nil {
		return j(err, "add more message failed")
	}

	err = tx.commit()
	if err != nil {
		return j(err, "commit failed")
	}

	return nil
}
