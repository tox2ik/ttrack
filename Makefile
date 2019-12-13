.PHONY: tt

all: tt install

tt:
	go build -ldflags="-s -w"

install:
	install -v ttrack ~/bin/tt
	upx --brute ~/bin/tt
