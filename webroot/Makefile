.PHONY: yarn-install dev build build-with-docker clean

yarn-install:
	@echo Getting dependencies using yarn
	yarn install

dev: yarn-install
	@echo Running webapp for development
	yarn run dev

build: yarn-install
	@echo Building webapp
	yarn run build

build-with-docker:
	@echo Building webapp with Docker
	docker run --rm -v $(PWD):/smsender -w /smsender node:6.10.0 sh -c "make build"

clean:
	@echo Cleaning webapp
	rm -f public/bundle.js
	rm -rf node_modules