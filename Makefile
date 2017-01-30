.PHONY: all fmt gofmt golint gotest install install-fcgi lint test	\
	vnu webapp-lint webapp-fmt webapp-jsdoc webapp-karma	\
	webapp-install webapp-webpack-watch

all: $(GOPATH)/bin/tsubonesystem3
	'$<' & $(MAKE) -C webapp webpack-watch & wait

fmt: webapp-fmt gofmt

gofmt:
	gofmt -w .

$(GOPATH)/bin/%: ./frontend/%
	go install $<

golint:
	golint ./...

gotest:
	go test ./...

install: webapp-install $(GOPATH)/bin/tsubonesystem3

install-fcgi: webapp-install $(GOPATH)/bin/tsubonesystem3_fcgi

lint: webapp-lint golint

prepare:
	$(MAKE) -C webapp prepare

test: webapp-karma gotest

webapp-vnu: $(GOPATH)/bin/tsubonesystem3 webapp-install
	'$<' & TSUBONESYSTEM_URL=`go run ./frontend/tsubonesystem3_resolve` $(MAKE) -C webapp vnu; kill $$!

webapp-lint: webapp
	$(MAKE) -C $< lint

webapp-fmt: webapp
	$(MAKE) -C $< fmt

webapp-karma:
	TSUBONESYSTEM_URL=`go run ./frontend/tsubonesystem3_resolve` $(MAKE) -C webapp karma

webapp-install: webapp
	$(MAKE) -C $< install

webapp-webpack-watch: webapp
	$(MAKE) -C $< webpack-watch
