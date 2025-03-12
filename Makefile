.PHONY: test run

test:
	go test -v -cover ./...

run:
	cd cmd && go run main.go

swagbuild:
	swag init -g cmd/main.go -o ./docs --parseDependency