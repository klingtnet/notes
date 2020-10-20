.PHONY: notes run embeds.go
LISTEN_ADDR?=localhost:3333

GIT_REVISION:=$(shell git describe --always --tags)

all: test notes

notes: embeds.go
	go build -ldflags="-X 'main.Version=$(GIT_REVISION)'" .

test: lint
	go test .

lint:
	golangci-lint run .

embeds.go:
	go run github.com/klingtnet/embed/cmd/embed --include assets --include migrations --include views

install: notes
	@install --strip -Dm700 notes $$HOME/.local/bin
	@install -Dm700 -d $$HOME/.config/notes
	@install -Dm600 dist/notes.systemd.service $$HOME/.config/systemd/user/notes.service
	@systemctl --user daemon-reload
	@echo Please copy dist/notes.env to $$HOME/.config/notes/notes.env and enter your database password
	@echo To start the service run: systemctl --user start notes.service
	@echo http://localhost:13333

run: notes
	./notes run --database-passphrase $$DATABASE_PASSPHRASE --database-path=notes.db --listen-addr=$(LISTEN_ADDR)

renew:
	./notes renew --database-passphrase $$DATABASE_PASSPHRASE --database-path=notes.db

rerun:
	git ls-files --cached | grep -v embeds.go | entr -c -r make run

clean:
	git clean -fd
