GO=go
GOCOVER=$(GO) tool cover
GOTEST=TZ=UTC $(GO) test

deps:
	go get -u ./...

vet:
	go vet ./...

test: vet
	$(GOTEST) ./... -cover -coverprofile=coverage.out

coverage:
	$(GOCOVER) -func=coverage.out
	@unlink coverage.out

lint:
	golangci-lint run -v

.PHONY: test/cover
test/cover:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCOVER) -func=coverage.out
	@unlink coverage.out

cover:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCOVER) -html=coverage.out -o coverage


development: dev

dev:
	air -c .air.toml

ci: test lint

update-golly-plugins:
	go get github.com/golly-go/plugins/eventsource
	go get github.com/golly-go/plugins/mongo
	go get github.com/golly-go/plugins/kafka

