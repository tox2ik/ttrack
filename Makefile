.PHONY: tt
.ONESHELL: tt test_mains test
PREFIX := ~

TRUN := go test

all: tt install
tt:; cd cmd/tt; go build -ldflags="-s -w"
install:
	install -v cmd/tt/tt $(PREFIX)/bin/tt
	#upx -qqq --lzma $(PREFIX)/bin/tt || true
testr:; make test TRUN='richgo test'
test: test_mains
	export tt_yes=1
	$(TRUN) ./...
cover:; $(TRUN) -cover ./... ; bash bin/code-percent.sh

test_mains:
	@export tt_yes=1
	for i in `grep -r . -e '^package main' -l | xargs -n1 dirname | sort -u`; do
	echo; cd $$i; echo IN: $$i;
	$(TRUN)
	echo
	done

