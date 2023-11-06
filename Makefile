.PHONY: all clean

binary   = srm
version  = 0.1.7
build	   = $(shell git rev-parse HEAD)
ldflags  = -ldflags "-X 'github.com/waldirborbajr/srm/command.version=$(version)'
ldflags += -X 'github.com/waldirborbajr/srm/command.build=$(build)'"

all:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/$(binary) $(ldflags) cmd/main.go

test:
	go test ./... -cover -coverprofile c.out
	go tool cover -html=c.out -o coverage.html

clean:
	rm -rf bin/$(binary) c.out coverage.html
