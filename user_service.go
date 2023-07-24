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
	magicKey, err := userService.randomIdService.generateRandomId()
	if err != nil {
		return nil, j(err, "generate magicKey failed")
	}
	accessToken, err := userService.randomIdService.generateRandomId()
	if err != nil {
		return nil, j(err, "generate accessToken failed")
	}

	tx, err := userService.dbService.beginTx(ctx)
	defer tx.rollBack()
	if err != nil {
		return nil, j(err, "start tx failed")
	}

	userId, err := tx.createUser(*magicKey)
	if err != nil {
		return nil, j(err, "create user failed")
	}

	err = tx.createAccessToken(*userId, *accessToken)
	if err != nil {
		return nil, j(err, "create accessToken failed")
	}

	err = tx.commit()
	if err != nil {
		return nil, j(err, "commit failed")
	}

	return &CreateUserSuccess{
		accessToken: *accessToken,
		magicKey:    *magicKey,
	}, nil
}
