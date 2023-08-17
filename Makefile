.SILENT: run

default: run

run:
	go run . ./dev-configuration.json

build_prod:
	go build -ldflags "-s -w"

gen:
	go run github.com/isaac-weisberg/go-jason-gen@v0.2.3 .

ffjson:
	go run github.com/pquerna/ffjson@v0.0.0-20190930134022-aa0246cd15f7 ./requests/access_token_having_legacy.go
