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
	./notes run --database-passphrase $$DATABASE_PASSPHRASE

renew:
	./notes renew --database-passphrase $$DATABASE_PASSPHRASE

rerun:
	git ls-files --cached | grep -v embeds.go | entr -c -r make run

clean:
	git clean -fd
