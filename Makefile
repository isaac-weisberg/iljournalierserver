.SILENT: run

default: run

run:
	go run . ./dev-configuration.json

build_prod:
	go build -ldflags "-s -w"

gen:
	go run github.com/isaac-weisberg/go-jason-gen@v0.2.3 .

easyjson:
	easyjson -all access_token_having_legacy.go	
	
# go run github.com/mailru/easyjson@v0.7.7 -all ./access_token_having_legacy.go
