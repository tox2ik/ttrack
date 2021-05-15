.PHONY: tt
.ONESHELL: tt
PREFIX := ~

all: tt install
tt:; cd cmd/tt; go build -ldflags="-s -w"
install:
	install -v cmd/tt/tt $(PREFIX)/bin/tt
	upx -qqq --lzma $(PREFIX)/bin/tt || true
testr:; richgo test ./...
test:; go test ./...
