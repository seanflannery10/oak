## audit: run quality control checks
.PHONY: audit
audit:
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go test -race -vet=off ./...
	go mod verify
	go mod vendor

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v