.PHONY: all clean

binary   = srm 
version  = 0.1.0
build	   = $(shell git rev-parse HEAD)
ldflags  = -ldflags "-X 'github.com/waldirborbajr/srm/command.version=$(version)'
ldflags += -X 'github.com/waldirborbajr/srm/command.build=$(build)'"

all:
	go build -o $(binary) $(ldflags) ./cmd/main.go
	mv $(binary) $(GOPATH)/bin
test:
	go test ./... -cover -coverprofile c.out
	go tool cover -html=c.out -o coverage.html

clean:
	rm -rf $(binary) c.out coverage.html
