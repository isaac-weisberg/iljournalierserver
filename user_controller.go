package main

import (
	"context"

	"caroline-weisberg.fun/iljournalierserver/errors"
	"caroline-weisberg.fun/iljournalierserver/services"
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
	}

	return &createUserResBody, nil
}

type loginRequestBody struct {
	LoginKey string `json:"loginKey" validate:"required"`
}

type loginResponseBody struct {
	accessTokenHavingLegacy
}

func (uc *userController) login(ctx context.Context, loginRequestBody *loginRequestBody) (*loginResponseBody, error) {
	loginSuccess, err := uc.userService.Login(loginRequestBody.LoginKey, ctx)
	if err != nil {
		return nil, errors.J(err, "user service login failed")
	}

	response := loginResponseBody{
		accessTokenHavingLegacy{
			AccessToken: loginSuccess.AccessToken,
		},
	}

	return &response, nil
}
