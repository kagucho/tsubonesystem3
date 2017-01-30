.PHONY: all fmt gofmt goinstall goinstall-fcgi golint gotest	\
	install install-fcgi lint test vnu webapp-lint webapp-fmt	\
	webapp-jsdoc webapp-karma webapp-install webapp-webpack-watch

all:
	"$$GOPATH/bin/tsubonesystem3" & $(MAKE) -C webapp webpack-watch & wait

fmt: webapp-fmt gofmt

gofmt:
	gofmt -w .

goinstall: frontend/tsubonesystem3
	go install ./$<

goinstall-fcgi: frontend/tsubonesystem3_fcgi
	go install ./$<

golint:
	golint ./...

gotest:
	go test ./...

install: webapp-install goinstall

install-fcgi: webapp-install goinstall-fcgi

lint: webapp-lint golint

prepare: webapp
	$(MAKE) -C $< prepare

test: webapp-karma gotest

webapp-vnu: goinstall webapp-install
	"$$GOPATH/bin/tsubonesystem3" & TSUBONESYSTEM_URL=`go run ./frontend/tsubonesystem3_resolve` $(MAKE) -C webapp vnu; kill $$!

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
