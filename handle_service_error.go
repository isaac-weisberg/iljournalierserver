package main

import (
	"net/http"

	"caroline-weisberg.fun/iljournalierserver/errors"
)

func handleServiceError(err error, w http.ResponseWriter, r *http.Request) {
	if errors.Is(err, errors.UserNotFoundForAccessToken) {
		w.WriteHeader(418)
		return
	} else if errors.Is(err, errors.UserNotFoundForMagicKey) {
		w.WriteHeader(418)
		return
	} else if errors.Is(err, errors.FlagDoesntBelongToTheUser) {
		w.WriteHeader(418)
		return
	} else {
		w.WriteHeader(500)
		return
	}
}
