package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type moreMessagesController struct {
	moreMessagesService *moreMessagesService
}

func newMoreMessagesController(moreMessagesService *moreMessagesService) moreMessagesController {
	return moreMessagesController{moreMessagesService: moreMessagesService}
}

type addMoreMessageRequestBody struct {
	accessTokenHavingObject
	Msg string `json:"msg"`
}

func (moreMessagesController *moreMessagesController) addMoreMessage(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	var addMoreMsgBody addMoreMessageRequestBody
	err = json.Unmarshal(body, &addMoreMsgBody)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	err = moreMessagesController.moreMessagesService.addMessage(r.Context(), addMoreMsgBody.AccessToken, addMoreMsgBody.Msg)

	if err != nil {
		if is(err, userNotFoundForAccessToken) {
			w.WriteHeader(418)
			return
		} else {
			w.WriteHeader(500)
			return
		}
	}

	w.WriteHeader(200)
	return
}
