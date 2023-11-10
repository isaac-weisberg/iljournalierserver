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
	PublicId    string
	Iv          string
}

func (userService *UserService) CreateUser(ctx context.Context) (*CreateUserSuccess, error) {
	publicUserId, err := userService.generatePublicId()
	if err != nil {
		return nil, errors.J(err, "generate publicId failed")
	}

	magicKey, err := userService.generateMagicKey()
	if err != nil {
		return nil, errors.J(err, "generate magicKey failed")
	}

	accessToken, err := userService.generateAccessToken()
	if err != nil {
		return nil, errors.J(err, "generate accessToken failed")
	}

	iv, err := userService.generateIv()
	if err != nil {
		return nil, errors.J(err, "generate iv failed")
	}

	return BeginTxBlock[CreateUserSuccess](userService.dbService, ctx, func(tx *transaction.Transaction) (*CreateUserSuccess, error) {
		if err != nil {
			return nil, errors.J(err, "start tx failed")
		}

		userId, err := tx.CreateUser(*publicUserId, *magicKey, *iv)
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
			PublicId:    *publicUserId,
			Iv:          *iv,
		}, nil
	})
}

type LoginSuccess struct {
	AccessToken string
	PublicId    string
	Iv          string
}

func (userService *UserService) Login(magicKey string, ctx context.Context) (*LoginSuccess, error) {
	return BeginTxBlock[LoginSuccess](userService.dbService, ctx, func(tx *transaction.Transaction) (*LoginSuccess, error) {
		user, err := tx.FindUserForMagicKey(magicKey)
		if err != nil {
			return nil, errors.J(err, "find user for magicKey failed")
		}

		if user == nil {
			return nil, errors.UserNotFoundForMagicKey
		}

		accessToken, err := userService.generateAccessToken()
		if err != nil {
			return nil, errors.J(err, "generate accessToken failed")
		}

		err = tx.CreateAccessToken(user.Id, *accessToken)
		if err != nil {
			return nil, errors.J(err, "creating access token entry failed")
		}

		return &LoginSuccess{
			AccessToken: *accessToken,
			PublicId:    user.PublicId,
			Iv:          user.Iv,
		}, nil
	})
}

func (userService *UserService) generateAccessToken() (*string, error) {
	accessToken, err := userService.randomIdService.GenerateRandomId()
	if err != nil {
		return nil, errors.J(err, "generateRandomId failed")
	}

	return accessToken, nil
}

func (userService *UserService) generateIv() (*string, error) {
	iv, err := userService.randomIdService.GenerateRandomId()
	if err != nil {
		return nil, errors.J(err, "generateIv failed")
	}

	return iv, nil
}

func (userService *UserService) generateMagicKey() (*string, error) {
	magicKey, err := userService.randomIdService.GenerateRandomId()
	if err != nil {
		return nil, errors.J(err, "generateRandomId failed")
	}

	return magicKey, nil
}

func (userService *UserService) generatePublicId() (*string, error) {
	publicId, err := userService.randomIdService.GenerateRandomId()
	if err != nil {
		return nil, errors.J(err, "generateRandomId failed")
	}

	return publicId, nil
}
