.PHONY: notes run embeds.go

GIT_REVISION:=$(shell git describe --always --tags)

all: notes lint

notes: embeds.go
	go build -ldflags="-X 'main.Version=$(GIT_REVISION)'" .

lint:
	golangci-lint run .

embeds.go:
	go run github.com/klingtnet/embed/cmd/embed --include assets --include migrations --include views

run: notes
	go run . --database-passphrase $$DATABASE_PASSPHRASE

clean:
	git clean -fd
