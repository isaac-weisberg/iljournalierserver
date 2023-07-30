.SILENT: run

default: run

run:
	go run . ./dev-configuration.json

build_prod:
	go build -ldflags "-s -w"