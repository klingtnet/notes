.PHONY: notes run embeds.go

GIT_REVISION:=$(shell git describe --always --tags)

all: notes

notes: embeds.go
	go build -ldflags="-X 'main.Version=$(GIT_REVISION)'" .

embeds.go:
	go run github.com/klingtnet/embed/cmd/embed --include assets --include migrations --include views

run: notes
	@./scripts/run

clean:
	git clean -fd
