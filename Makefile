

.PHONY: clean all init generate generate_mocks

all: build/main

run: 
	DATABASE_URL="postgres://postgres:postgres@localhost:5432/database?sslmode=disable" go run cmd/main.go

seed:
	DATABASE_URL="postgres://postgres:postgres@localhost:5432/database?sslmode=disable" go run cmd/main.go seed

build/main: cmd/main.go generated
	@echo "Building..."
	go build -o $@ $<

clean:
	rm -rf generated

init: clean generate
	go mod tidy
	go mod vendor

test:
	go clean -testcache
	go test -short -coverprofile coverage.out -short -v ./...

test_api:
	go clean -testcache
	go test ./tests/...

generate: generated generate_mocks generate_mocks

generated: api.yml
	@echo "Generating files..."
	mkdir generated || true
	oapi-codegen --package generated -generate types,server,spec $< > generated/api.gen.go

generate_mocks: 
	go generate ./...