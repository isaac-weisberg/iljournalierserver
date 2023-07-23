package main

import (
	"encoding/json"
	"net/http"
)

type userController struct {
	userService *userService
}

func newUserController(userService *userService) userController {
	return userController{userService: userService}
}

type CreateUserResponseBody struct {
	LoginKey    string `json:"loginKey"`
	AccessToken string `json:"accessToken"`
}

func (uc *userController) createUser(w http.ResponseWriter, r *http.Request) {
	user, err := uc.userService.createUser(r.Context())

	if err != nil {
		w.WriteHeader(500)
		return
	}

	createUserResBody := CreateUserResponseBody{
		LoginKey:    user.magicKey,
		AccessToken: user.accessToken,
	}

	bytes, err := json.Marshal(createUserResBody)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	w.Write(bytes)
}
