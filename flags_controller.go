package main

import "net/http"

type flagsController struct {
	flagsService *flagsService
}

func newFlagsController(flagsService *flagsService) flagsController {
	return flagsController{
		flagsService: flagsService,
	}
}

type markFlagsRequestBodyMarkRequest struct {
	UnixSeconds int64 `json:"unixSeconds"`
	FlagId      int64 `json:"flagId"`
}

type markFlagsRequestBody struct {
	Requests []markFlagsRequestBodyMarkRequest `json:"requests"`
}

func (flagsController *flagsController) markFlags(w http.ResponseWriter, r *http.Request) {

}
