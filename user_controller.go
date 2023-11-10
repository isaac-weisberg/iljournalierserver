package main

import (
	"context"

	"caroline-weisberg.fun/iljournalierserver/errors"
	"caroline-weisberg.fun/iljournalierserver/services"
	gojason "github.com/isaac-weisberg/go-jason"
)

type userController struct {
	userService *services.UserService
}

func newUserController(userService *services.UserService) userController {
	return userController{userService: userService}
}

type createUserResponseBody struct {
	accessTokenHavingLegacy
	LoginKey string `json:"loginKey" validate:"required"`
	PublicId string `json:"publicId" validate:"required"`
	Iv       string `json:"iv" validate:"required"`
}

func (uc *userController) createUser(ctx context.Context) (*createUserResponseBody, error) {
	user, err := uc.userService.CreateUser(ctx)

	if err != nil {
		return nil, errors.J(err, "create user service failed")
	}

	createUserResBody := createUserResponseBody{
		accessTokenHavingLegacy: accessTokenHavingLegacy{
			AccessToken: user.AccessToken,
		},
		LoginKey: user.MagicKey,
		PublicId: user.PublicId,
		Iv:       user.Iv,
	}

	return &createUserResBody, nil
}

type loginRequestBody struct {
	gojason.Decodable

	loginKey string
}

type loginResponseBody struct {
	accessTokenHavingLegacy
	PublicId string `json:"publicId"`
	Iv       string `json:"iv"`
}

func (uc *userController) login(ctx context.Context, loginRequestBody *loginRequestBody) (*loginResponseBody, error) {
	loginSuccess, err := uc.userService.Login(loginRequestBody.loginKey, ctx)
	if err != nil {
		return nil, errors.J(err, "user service login failed")
	}

	response := loginResponseBody{
		accessTokenHavingLegacy{
			AccessToken: loginSuccess.AccessToken,
		},
		loginSuccess.PublicId,
		loginSuccess.Iv,
	}

	return &response, nil
}
