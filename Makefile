## audit: run quality control checks
.PHONY: audit
audit:
	go mod vendor
	go mod tidy -v
	golangci-lint run --fix
	go test -race -vet=off ./...
	go mod verify

## coverage: test and view code coverage
.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

## audit: run quality control checks
.PHONY: upgrade
upgrade:
	go get -u ./...