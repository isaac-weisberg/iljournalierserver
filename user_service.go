package main

import (
	"context"

	"caroline-weisberg.fun/iljournalierserver/errors"
	"caroline-weisberg.fun/iljournalierserver/transaction"
)

type userService struct {
	dbService       *databaseService
	randomIdService *randomIdService
}

func newUserService(dbService *databaseService, randomIdService *randomIdService) userService {
	return userService{dbService: dbService, randomIdService: randomIdService}
}

type createUserSuccess struct {
	accessToken string
	magicKey    string
}

func (userService *userService) createUser(ctx context.Context) (*createUserSuccess, error) {
	magicKey, err := userService.generateMagicKey()
	if err != nil {
		return nil, errors.J(err, "generate magicKey failed")
	}

	accessToken, err := userService.generateAccessToken()
	if err != nil {
		return nil, errors.J(err, "generate accessToken failed")
	}

	return beginTxBlock[createUserSuccess](userService.dbService, ctx, func(tx *transaction.Transaction) (*createUserSuccess, error) {
		if err != nil {
			return nil, errors.J(err, "start tx failed")
		}

		userId, err := tx.CreateUser(*magicKey)
		if err != nil {
			return nil, errors.J(err, "create user failed")
		}

		err = tx.CreateAccessToken(*userId, *accessToken)
		if err != nil {
			return nil, errors.J(err, "create accessToken failed")
		}

		return &createUserSuccess{
			accessToken: *accessToken,
			magicKey:    *magicKey,
		}, nil
	})
}

type loginSuccess struct {
	accessToken string
}

func (userService *userService) login(magicKey string, ctx context.Context) (*loginSuccess, error) {
	return beginTxBlock[loginSuccess](userService.dbService, ctx, func(tx *transaction.Transaction) (*loginSuccess, error) {
		userId, err := tx.FindUserForMagicKey(magicKey)
		if err != nil {
			return nil, errors.J(err, "find user for magicKey failed")
		}

		if userId == nil {
			return nil, errors.UserNotFoundForMagicKey
		}

		accessToken, err := userService.generateAccessToken()
		if err != nil {
			return nil, errors.J(err, "generate accessToken failed")
		}

		err = tx.CreateAccessToken(*userId, *accessToken)
		if err != nil {
			return nil, errors.J(err, "creating access token entry failed")
		}

		return &loginSuccess{accessToken: *accessToken}, nil
	})
}

func (userService *userService) generateAccessToken() (*string, error) {
	accessToken, err := userService.randomIdService.generateRandomId()
	if err != nil {
		return nil, errors.J(err, "generateRandomId failed")
	}

	return accessToken, nil
}

func (userService *userService) generateMagicKey() (*string, error) {
	magicKey, err := userService.randomIdService.generateRandomId()
	if err != nil {
		return nil, errors.J(err, "generateRandomId failed")
	}

	return magicKey, nil
}
