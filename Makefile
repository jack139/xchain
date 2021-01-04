BUILD=build
BUILD_TAGS=cleveldb
CGO_ENABLED=1

all: clean build

build:
	CGO_ENABLED=$(CGO_ENABLED) go build -o $(BUILD)/xchain -tags '$(BUILD_TAGS)' xchain/main.go
	CGO_ENABLED=$(CGO_ENABLED) go build -o $(BUILD)/xcli  xcli/main.go

clean:
	@rm -rf $(BUILD)
	@go clean
