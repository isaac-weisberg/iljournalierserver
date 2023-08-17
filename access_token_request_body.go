package main

import gojason "github.com/isaac-weisberg/go-jason"

type accessTokenHavingRequest struct {
	gojason.Decodable

	accessToken string
}
