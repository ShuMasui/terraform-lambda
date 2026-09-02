.PHONY: deps build test zip deploy

deps:
	go mod download
	go mod tidy

build:
	make deps
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o bin/bootstrap cmd/lambda/main.go

test:
	go test ./...

zip:
	make build
	zip -j bin/lambda.zip bin/bootstrap

deploy:
	make zip
	aws --profile study-aws lambda update-function-code --region us-east-1 --function-name api --zip-file fileb://bin/lambda.zip