.PHONY: yarninstall dev build build-with-docker clean

yarninstall:
	@echo Getting dependencies using yarn
	yarn install

dev: yarninstall
	@echo Running webapp for development
	yarn run dev

build: yarninstall
	@echo Building webapp
	yarn run build

build-with-docker:
	@echo Building webapp with Docker
	docker run --rm -v $(PWD):/smsender -w /smsender node:6.10.0 sh -c "make build"

clean:
	@echo Cleaning webapp
	rm -f public/bundle.js
	rm -rf node_modules