.PHONY: tt

tt:
	go build -ldflags="-s -w"

install:
	install -v tt ~/bin/tt
	upx --brute ~/bin/tt
