package main

import (
	"context"

	"caroline-weisberg.fun/iljournalierserver/errors"
)

type flagsService struct {
	databaseService *databaseService
}

func newFlagsService(databaseService *databaseService) flagsService {
	return flagsService{
		databaseService: databaseService,
	}
}

func (flagsService *flagsService) markFlags(ctx context.Context, accessToken string, markFlagRequests []markFlagRequest) error {
	return beginTxBlockVoid(flagsService.databaseService, ctx, func(tx *transaction) error {
		userId, err := tx.findUserIdForAccessToken(accessToken)
		if err != nil {
			return errors.J(err, "findUserIdForAccessToken failed")
		}

		if userId == nil {
			return errors.UserNotFoundForAccessToken
		}

		flagIds, err := tx.getKnownFlagIdsForUser(*userId)
		if err != nil {
			return errors.J(err, "getKnownFlagIdsForUser failed")
		}

		var flagIdsMap = mapFromSlice[int64](flagIds)

		for _, markedFlag := range markFlagRequests {
			if !mapContains[int64](flagIdsMap, markedFlag.FlagId) {
				return errors.FlagDoesntBelongToTheUser
			}
		}

		err = tx.markFlags(markFlagRequests)
		if err != nil {
			return errors.J(err, "failed marking flags")
		}

		return nil
	})
}
