.PHONY: all fmt gofmt goinstall goinstall-fcgi golint gotest	\
	install install-fcgi lint test vnu webapp-eslint	\
	webapp-eslint-fix webapp-jsdoc webapp-karma	\
	webapp-webpack webapp-webpack-watch

all:
	"$$GOPATH/bin/tsubonesystem3" & $(MAKE) -C webapp webpack-watch & wait

fmt: webapp-eslint-fix gofmt

gofmt:
	gofmt -w .

goinstall:
	go install ./frontend/tsubonesystem3

goinstall-fcgi:
	go install ./frontend/tsubonesystem3_fcgi

golint:
	golint ./...

gotest:
	go test ./...

install: webapp-webpack goinstall

install-fcgi: webapp-webpack goinstall-fcgi

lint: webapp-eslint vnu golint

test: webapp-karma gotest

vnu:
	$(MAKE) start
	TSUBONESYSTEM_URL=`go run ./frontend/tsubonesystem3_resolve` java -jar node_modules/vnu-jar/build/dist/vnu.jar "$$TSUBONESYSTEM_URL" "$$TSUBONESYSTEM_URL/license" "$$TSUBONESYSTEM_URL/private" "$$TSUBONESYSTEM_URL/private?_escaped_fragment_="
	$(MAKE) stop

webapp-eslint:
	$(MAKE) -C webapp eslint

webapp-eslint-fix:
	$(MAKE) -C webapp eslint-fix

webapp-karma:
	TSUBONESYSTEM_URL=`go run ./frontend/tsubonesystem3_resolve` $(MAKE) -C webapp karma

webapp-webpack:
	$(MAKE) -C webapp webpack

webapp-webpack-watch:
	$(MAKE) -C webapp webpack-watch
