package services

import (
	"context"

	"caroline-weisberg.fun/iljournalierserver/errors"
	"caroline-weisberg.fun/iljournalierserver/transaction"
)

type MoreMessagesService struct {
	databaseService *DatabaseService
}

func NewMoreMessagesService(databaseService *DatabaseService) MoreMessagesService {
	return MoreMessagesService{databaseService: databaseService}
}

func (moreMessagesService *MoreMessagesService) AddMessage(ctx context.Context, accessToken string, unixSeconds int64, msg string) error {
	return BeginTxBlockVoid(moreMessagesService.databaseService, ctx, func(tx *transaction.Transaction) error {
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
