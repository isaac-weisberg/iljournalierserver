package main

type accessTokenHavingObject struct {
	AccessToken string `json:"accessToken" validate:"required"`
}
