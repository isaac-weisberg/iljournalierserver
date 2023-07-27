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

func (moreMessagesService *moreMessagesService) addMessage(ctx context.Context, accessToken string, unixSeconds int64, msg string) error {
	return beginTxBlockVoid(moreMessagesService.databaseService, ctx, func(tx *transaction) error {
		userId, err := tx.findUserIdForAccessToken(accessToken)
		if err != nil {
			return j(err, "find user for accessToken failed")
		}

		if userId == nil {
			return userNotFoundForAccessToken
		}

		err = tx.addMoreMessage(*userId, unixSeconds, msg)
		if err != nil {
			return j(err, "add more message failed")
		}

		return nil
	})
}
