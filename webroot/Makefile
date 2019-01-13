.PHONY: yarn-install lint dev build build-with-docker clean

yarn-install:
	@echo Getting dependencies using yarn
	yarn install

lint:
	@echo Checking lint
	yarn run lint

dev: yarn-install
	@echo Running webapp for development
	yarn run dev

build: yarn-install
	@echo Building webapp
	yarn run build

build-with-docker: clean
	@echo Building webapp with Docker
	docker run --rm -v $(PWD):/smsender -w /smsender node:10.15.0 sh -c "make build"

clean:
	@echo Cleaning webapp
	rm -rf ./dist ./node_modules