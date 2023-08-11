package main

import gojason "github.com/isaac-weisberg/go-jason"

type accessTokenHavingRequest struct {
	gojason.Decodable

	accessToken string
}

type accessTokenHavingLegacy struct {
	AccessToken string `json:"accessToken" validate:"required"`
}
