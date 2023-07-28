package main

import (
	"context"

	"caroline-weisberg.fun/iljournalierserver/errors"
	"caroline-weisberg.fun/iljournalierserver/services"
	"caroline-weisberg.fun/iljournalierserver/transaction"
)

type flagsService struct {
	databaseService *services.DatabaseService
}

func newFlagsService(databaseService *services.DatabaseService) flagsService {
	return flagsService{
		databaseService: databaseService,
	}
}

func (flagsService *flagsService) markFlags(ctx context.Context, accessToken string, markFlagRequests []transaction.MarkFlagRequest) error {
	return services.BeginTxBlockVoid(flagsService.databaseService, ctx, func(tx *transaction.Transaction) error {
		userId, err := tx.FindUserIdForAccessToken(accessToken)
		if err != nil {
			return errors.J(err, "findUserIdForAccessToken failed")
		}

		if userId == nil {
			return errors.UserNotFoundForAccessToken
		}

		flagIds, err := tx.GetKnownFlagIdsForUser(*userId)
		if err != nil {
			return errors.J(err, "getKnownFlagIdsForUser failed")
		}

		var flagIdsMap = mapFromSlice[int64](flagIds)

		for _, markedFlag := range markFlagRequests {
			if !mapContains[int64](flagIdsMap, markedFlag.FlagId) {
				return errors.FlagDoesntBelongToTheUser
			}
		}

		err = tx.MarkFlags(markFlagRequests)
		if err != nil {
			return errors.J(err, "failed marking flags")
		}

		return nil
	})
}
