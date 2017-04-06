.PHONY: all doc eslint eslint-build eslint-fix eslint-fix-build	\
	eslint-fix-src eslint-fix-test eslint-src eslint-test fmt	\
	install jsdoc karma lint prepare webpack webpack-watch

all: eslint-test

doc: jsdoc

eslint: eslint-build eslint-src eslint-test

eslint-build:
	npm run eslint-build

eslint-fix: eslint-fix-build eslint-fix-src eslint-fix-test

eslint-fix-build:
	npm run eslint-build -- --fix

eslint-fix-src:
	npm run eslint-src -- --fix

eslint-fix-test:
	npm run eslint-test -- --fix

eslint-src:
	npm run eslint-src

eslint-test:
	npm run eslint-test

fmt: eslint-fix

install: webpack

jsdoc:
	npm run jsdoc

karma:
	npm run karma -- start --browsers PhantomJS --single-run

lint: eslint vnu

mithril:
	cd node_modules/mithril; npm run build

prepare:
	npm update --dev

test: karma

vnu:
	npm run vnu "$$TSUBONESYSTEM_URL" "$$TSUBONESYSTEM_URL/license" "$$TSUBONESYSTEM_URL/private" "$$TSUBONESYSTEM_URL/private?_escaped_fragment_="

webpack: mithril
	npm run webpack -- --output-path "$$GOPATH/share/tsubonesystem3/public" $(WEBPACKFLAGS)

webpack-watch: mithril
	npm run webpack -- --output-path "$$GOPATH/share/tsubonesystem3/public" --watch $(WEBPACKFLAGS)
