## audit: run quality control checks
.PHONY: audit
audit:
	go mod vendor
	go mod tidy -v
	go fmt ./...
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go test -race -vet=off ./...
	go mod verify

## coverage: test and view code coverage
.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out