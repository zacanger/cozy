PREFIX ?= /usr/local
CONFIG_PREFIX ?= /usr/share

cozy:
	go build

install:
	mkdir -p $(PREFIX)/bin
	cp -f cozy $(PREFIX)/bin/cozy
	chmod 755 $(PREFIX)/bin/cozy

clean:
	rm -f cozy

test:
	go test ./...

count:
	cloc --exclude-dir=x --force-lang=Ruby,cz . | sed 's/Ruby/cozy/g'

.PHONY: cozy test clean install count