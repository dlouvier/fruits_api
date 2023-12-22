# Define variables
GOCMD := go
GOBUILD := $(GOCMD) build
GOBUILDFLAGS := CGO_ENABLED=0
GOTEST := $(GOCMD) test
GORUN := $(GOCMD) run app.go
BINARY_NAME := ../fruits_api
SWAGCMD := /home/dani/go/bin/swag
SWAGINIT := $(SWAGCMD) init -g app.go
SWAGFMT := $(SWAGCMD) init -g app.go
TEST_FLAGS := -v -cover

# Default target when running `make` without any arguments
default: run

# Build the project and compile it in a single binary
build:
	@cd ./src && goreleaser --snapshot --clean

# Run tests
test:
	@cd ./src && $(GOTEST) $(TEST_FLAGS) ./...
	@cd ./src && goreleaser check

# Run the webserver (without building)
run:
	@trap 'echo "SIGTERM received, stopping..."; kill $$!; exit' INT; \
	cd ./src && $(GORUN) & wait

# Generate docs & prepare swagger UI
generate-docs:
	@cd ./src && $(SWAGINIT) && $(SWAGFMT)
