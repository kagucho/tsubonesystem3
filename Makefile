SELINUX = /usr/share/selinux/devel/Makefile

.PHONY: all fmt gofmt goinstall goinstall-fcgi golint gotest install	\
	install-fcgi lint test vnu webapp-fmt webapp-install	\
	webapp-lint webapp-test

all:
	"$$GOPATH/bin/tsubonesystem3" & $(MAKE) -C webapp webpack-watch & wait

doc: godoc jsdoc

fix: gofix

fmt: webapp-fmt gofmt

godoc: jsdoc
	mkdir -p go/src/github.com/kagucho/tsubonesystem3
	git ls-files | xargs cp --parent -t go/src/github.com/kagucho/tsubonesystem3
	(export ADDRESS=`go run ./frontend/tsubonesystem3_resolve/*`; GOPATH= godoc -goroot go -http $$ADDRESS & sleep 1; wget -nH -kmP out/go $$ADDRESS; kill $$!)

godoc-clean:
	rm -rf go out

gofix:
	go fix ./...

gofmt:
	go fmt ./...

govet:
	go vet ./...

goinstall: frontend/tsubonesystem3
	go install ./$<

goinstall-fcgi: frontend/tsubonesystem3_fcgi
	go install ./$<

golint:
	golint ./...

goprepare:
	go get -d ./...

gotest:
	go test ./...

install: webapp-install goinstall

install-fcgi: webapp-install goinstall-fcgi

jsdoc:
	cd webapp; npm run jsdoc -- -d ../out $(JSDOCFLAGS)

lint: webapp-lint golint

prepare: goprepare webapp-prepare

test: webapp-test gotest

vet: govet

webapp-fmt: webapp
	$(MAKE) -C $< fmt

webapp-install: webapp
	$(MAKE) -C $< install

webapp-lint: webapp
	"$$GOPATH/bin/tsubonesystem3" & $(MAKE) -C $< lint TSUBONESYSTEM_URL=http://`go run ./frontend/tsubonesystem3_resolve/*`; kill $$!

webapp-prepare: webapp
	$(MAKE) -C $< prepare

webapp-test:
	"$$GOPATH/bin/tsubonesystem3" & $(MAKE) -C webapp test TSUBONESYSTEM_URL=http://`go run ./frontend/tsubonesystem3_resolve/*`; kill $$!

%.pp: $(SELINUX)
	$(MAKE) -f $< $@
