package main

import gojason "github.com/isaac-weisberg/go-jason"

type accessTokenHavingObject struct {
	gojason.Decodable

	AccessToken string `json:"accessToken" validate:"required"`
}
