.PHONY: all eslint eslint-build eslint-fix eslint-fix-build	\
	eslint-fix-src eslint-fix-test eslint-src eslint-test fmt	\
	install jsdoc karma lint prepare webpack webpack-watch

all: eslint-test

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

prepare:
	npm update --dev

vnu:
	npm run vnu "$$TSUBONESYSTEM_URL" "$$TSUBONESYSTEM_URL/license" "$$TSUBONESYSTEM_URL/private" "$$TSUBONESYSTEM_URL/private?_escaped_fragment_="

webpack:
	npm run webpack -- --output-path "$$GOPATH/share/tsubonesystem3/public"

webpack-watch:
	npm run webpack -- --output-path "$$GOPATH/share/tsubonesystem3/public" --watch
