.PHONY: tt
.ONESHELL: tt test_mains test
PREFIX := ~

TRUN := go test

all: tt install ; echo 'alias tti="make && tt i"'
tt:; cd cmd/tt; go build -ldflags="-s -w"
install:
	install -v cmd/tt/tt $(PREFIX)/bin/tt
	#upx -qqq --lzma $(PREFIX)/bin/tt || true
rest:; make test TRUN='richgo test'
test: test_mains ; @$(TRUN) ./...
cover:; $(TRUN) -cover ./... ; bash bin/code-percent.sh
test_mains:
	@for i in `find -name main.go | xargs -n1 dirname | sort -u`; do
	echo; cd $$i; echo IN: $$i;
	$(TRUN) ./...
	echo
	cd $(PWD)
	done

clean:
	go clean -testcache
