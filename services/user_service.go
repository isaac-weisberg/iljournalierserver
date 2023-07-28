package services

import (
	"context"

	"caroline-weisberg.fun/iljournalierserver/errors"
	"caroline-weisberg.fun/iljournalierserver/transaction"
)

type UserService struct {
	dbService       *DatabaseService
	randomIdService *RandomIdService
}

func NewUserService(dbService *DatabaseService, randomIdService *RandomIdService) UserService {
	return UserService{dbService: dbService, randomIdService: randomIdService}
}

type CreateUserSuccess struct {
	AccessToken string
	MagicKey    string
}

func (userService *UserService) CreateUser(ctx context.Context) (*CreateUserSuccess, error) {
	magicKey, err := userService.generateMagicKey()
	if err != nil {
		return nil, errors.J(err, "generate magicKey failed")
	}

	accessToken, err := userService.generateAccessToken()
	if err != nil {
		return nil, errors.J(err, "generate accessToken failed")
	}

	return BeginTxBlock[CreateUserSuccess](userService.dbService, ctx, func(tx *transaction.Transaction) (*CreateUserSuccess, error) {
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

		return &CreateUserSuccess{
			AccessToken: *accessToken,
			MagicKey:    *magicKey,
		}, nil
	})
}

type LoginSuccess struct {
	AccessToken string
}

func (userService *UserService) Login(magicKey string, ctx context.Context) (*LoginSuccess, error) {
	return BeginTxBlock[LoginSuccess](userService.dbService, ctx, func(tx *transaction.Transaction) (*LoginSuccess, error) {
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

		return &LoginSuccess{AccessToken: *accessToken}, nil
	})
}

func (userService *UserService) generateAccessToken() (*string, error) {
	accessToken, err := userService.randomIdService.GenerateRandomId()
	if err != nil {
		return nil, errors.J(err, "generateRandomId failed")
	}

	return accessToken, nil
}

func (userService *UserService) generateMagicKey() (*string, error) {
	magicKey, err := userService.randomIdService.GenerateRandomId()
	if err != nil {
		return nil, errors.J(err, "generateRandomId failed")
	}

	return magicKey, nil
}
