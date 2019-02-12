GOPACKAGES = $(shell go list ./...)

default: run

run:
	@go run ./cmd/appname

setup:
	sh ./setup.sh

test: test-all

test-all:
	@go test $(GOPACKAGES) -coverprofile=coverage.out

test-ci:
	@go test -v -tags integration $(GOPACKAGES) -coverprofile=$(TEST_RESULTS)/coverage.out

view-coverage:
	@go tool cover -html=coverage.out

test-integration:
	@go test -tags integration $(GOPACKAGES)
