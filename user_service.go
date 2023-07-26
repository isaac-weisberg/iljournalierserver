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

type createUserSuccess struct {
	accessToken string
	magicKey    string
}

func (userService *userService) createUser(ctx context.Context) (*createUserSuccess, error) {
	magicKey, err := userService.generateMagicKey()
	if err != nil {
		return nil, j(err, "generate magicKey failed")
	}

	accessToken, err := userService.generateAccessToken()
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

	return &createUserSuccess{
		accessToken: *accessToken,
		magicKey:    *magicKey,
	}, nil
}

type loginSuccess struct {
	accessToken string
}

func (userService *userService) login(magicKey string, ctx context.Context) (*loginSuccess, error) {
	tx, err := userService.dbService.beginTx(ctx)
	if err != nil {
		return nil, j(err, "creating tx failed")
	}

	userId, err := tx.findUserForMagicKey(magicKey)
	if err != nil {
		return nil, j(err, "find user for magicKey failed")
	}

	if userId == nil {
		return nil, userNotFoundForMagicKey
	}

	accessToken, err := userService.generateAccessToken()
	if err != nil {
		return nil, j(err, "generate accessToken failed")
	}

	err = tx.createAccessToken(*userId, *accessToken)
	if err != nil {
		return nil, j(err, "creating access token entry failed")
	}

	err = tx.commit()
	if err != nil {
		return nil, j(err, "commit failed")
	}

	return &loginSuccess{accessToken: *accessToken}, nil
}

func (userService *userService) generateAccessToken() (*string, error) {
	accessToken, err := userService.randomIdService.generateRandomId()
	if err != nil {
		return nil, j(err, "generateRandomId failed")
	}

	return accessToken, nil
}

func (userService *userService) generateMagicKey() (*string, error) {
	magicKey, err := userService.randomIdService.generateRandomId()
	if err != nil {
		return nil, j(err, "generateRandomId failed")
	}

	return magicKey, nil
}
