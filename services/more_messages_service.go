package services

import (
	"context"

	"caroline-weisberg.fun/iljournalierserver/errors"
	"caroline-weisberg.fun/iljournalierserver/models"
	"caroline-weisberg.fun/iljournalierserver/transaction"
)

type MoreMessagesService struct {
	databaseService *DatabaseService
}

func NewMoreMessagesService(databaseService *DatabaseService) MoreMessagesService {
	return MoreMessagesService{databaseService: databaseService}
}

func (moreMessagesService *MoreMessagesService) AddMessage(
	ctx context.Context,
	accessToken string,
	addMessageRequests []models.AddMessageRequest,
) error {
	return BeginTxBlockVoid(moreMessagesService.databaseService, ctx, func(tx *transaction.Transaction) error {
		userId, err := tx.FindUserIdForAccessToken(accessToken)
		if err != nil {
			return errors.J(err, "find user for accessToken failed")
		}

		if userId == nil {
			return errors.UserNotFoundForAccessToken
		}

		err = tx.AddMoreMessages(*userId, addMessageRequests)
		if err != nil {
			return errors.J(err, "add more message failed")
		}

		return nil
	})
}
