.PHONY: all eslint eslint-build eslint-fix eslint-fix-build	\
	eslint-fix-src eslint-fix-test eslint-src eslint-test	\
	jsdoc karma webpack webpack-watch

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

jsdoc:
	npm run jsdoc

karma:
	npm run karma -- start --browsers PhantomJS --single-run

webpack:
	npm run webpack -- --output-path "$$GOPATH/share/tsubonesystem3/public"

webpack-watch:
	npm run webpack -- --output-path "$$GOPATH/share/tsubonesystem3/public" --watch
