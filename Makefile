.PHONY: tt

all: tt install

tt:
	go build -ldflags="-s -w"

install:
	install -v ttrack ~/bin/tt
	upx -qqq --lzma ~/bin/tt

pack:
	upx --brute ~/bin/tt
