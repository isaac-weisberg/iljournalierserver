package main

import (
	"encoding/json"
	"io"
	"net/http"

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
	accessTokenHavingObject
	LoginKey string `json:"loginKey"`
}

func (uc *userController) createUser(w http.ResponseWriter, r *http.Request) {
	user, err := uc.userService.CreateUser(r.Context())

	if err != nil {
		w.WriteHeader(500)
		return
	}

	createUserResBody := createUserResponseBody{
		accessTokenHavingObject: accessTokenHavingObject{
			AccessToken: user.AccessToken,
		},
		LoginKey: user.MagicKey,
	}

	bytes, err := json.Marshal(createUserResBody)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	w.Write(bytes)
}

type loginRequestBody struct {
	LoginKey string `json:"loginKey"`
}

type loginResponseBody struct {
	accessTokenHavingObject
}

func (uc *userController) login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	var loginRequestBody loginRequestBody
	err = json.Unmarshal(body, &loginRequestBody)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	loginSuccess, err := uc.userService.Login(loginRequestBody.LoginKey, r.Context())
	if err != nil {
		if errors.Is(err, errors.UserNotFoundForMagicKey) {
			w.WriteHeader(418)
			return
		} else {
			w.WriteHeader(500)
			return
		}
	}

	response := loginResponseBody{
		accessTokenHavingObject{
			AccessToken: loginSuccess.AccessToken,
		},
	}

	resData, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	w.Write(resData)
}
