package main

import (
	"context"

	"caroline-weisberg.fun/iljournalierserver/errors"
	"caroline-weisberg.fun/iljournalierserver/transaction"
)

type moreMessagesService struct {
	databaseService *databaseService
}

func newMoreMessagesService(databaseService *databaseService) moreMessagesService {
	return moreMessagesService{databaseService: databaseService}
}

func (moreMessagesService *moreMessagesService) addMessage(ctx context.Context, accessToken string, unixSeconds int64, msg string) error {
	return beginTxBlockVoid(moreMessagesService.databaseService, ctx, func(tx *transaction.Transaction) error {
		userId, err := tx.FindUserIdForAccessToken(accessToken)
		if err != nil {
			return errors.J(err, "find user for accessToken failed")
		}

		if userId == nil {
			return errors.UserNotFoundForAccessToken
		}

		err = tx.AddMoreMessage(*userId, unixSeconds, msg)
		if err != nil {
			return errors.J(err, "add more message failed")
		}

		return nil
	})
}
