package services

import (
	"context"

	"caroline-weisberg.fun/iljournalierserver/errors"
	"caroline-weisberg.fun/iljournalierserver/transaction"
	"caroline-weisberg.fun/iljournalierserver/utils"
)

type FlagsService struct {
	databaseService *DatabaseService
}

func NewFlagsService(databaseService *DatabaseService) FlagsService {
	return FlagsService{
		databaseService: databaseService,
	}
}

func (flagsService *FlagsService) MarkFlags(ctx context.Context, accessToken string, markFlagRequests []transaction.MarkFlagRequest) error {
	return BeginTxBlockVoid(flagsService.databaseService, ctx, func(tx *transaction.Transaction) error {
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

		var flagIdsMap = utils.MapFromSlice[int64](flagIds)

		for _, markedFlag := range markFlagRequests {
			if !utils.MapContains[int64](flagIdsMap, markedFlag.FlagId) {
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

func (flagsService *FlagsService) AddKnownFlags(ctx context.Context, accessToken string, newFlags []string) (*[]int64, error) {
	return BeginTxBlock[[]int64](flagsService.databaseService, ctx, func(tx *transaction.Transaction) (*[]int64, error) {
		userId, err := tx.FindUserIdForAccessToken(accessToken)
		if err != nil {
			return nil, errors.J(err, "find user for access token failed")
		}

		if userId == nil {
			return nil, errors.UserNotFoundForAccessToken
		}

		flagIds, err := tx.AddKnownFlags(*userId, newFlags)
		if err != nil {
			return nil, errors.J(err, "tx add known flags failed")
		}

		return flagIds, nil
	})
}
