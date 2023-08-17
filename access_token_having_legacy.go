package main

type accessTokenHavingLegacy struct {
	AccessToken string `json:"accessToken" validate:"required"`
}
