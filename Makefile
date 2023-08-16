.SILENT: run

default: run

run:
	go run . ./dev-configuration.json

build_prod:
	go build -ldflags "-s -w"

gen:
	go run github.com/isaac-weisberg/go-jason-gen@v0.2.2 .
