package main

import (
	"context"
)

type userService struct {
	dbService       *databaseService
	randomIdService *randomIdService
}

func newUserService(dbService *databaseService, randomIdService *randomIdService) userService {
	return userService{dbService: dbService, randomIdService: randomIdService}
}

type CreateUserSuccess struct {
	accessToken string
	magicKey    string
}

func (userService *userService) createUser(ctx context.Context) (*CreateUserSuccess, error) {
	wrapError := createErrorWrapper("userServiceCreateUserError")

	magicKey, err := userService.randomIdService.generateRandomId()
	if err != nil {
		return nil, wrapError(err)
	}
	accessToken, err := userService.randomIdService.generateRandomId()
	if err != nil {
		return nil, wrapError(err)
	}

	tx, err := userService.dbService.beginTx(ctx)
	defer tx.rollBack()
	if err != nil {
		return nil, wrapError(err)
	}

	userId, err := tx.createUser(*magicKey)
	if err != nil {
		return nil, wrapError(err)
	}

	err = tx.createAccessToken(*userId, *accessToken)
	if err != nil {
		return nil, wrapError(err)
	}

	err = tx.commit()
	if err != nil {
		return nil, wrapError(err)
	}

	return &CreateUserSuccess{accessToken: *accessToken, magicKey: *magicKey}, nil
}
