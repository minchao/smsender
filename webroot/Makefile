.PHONY: build docker-build dev clean

.npminstall: package.json
	@echo Getting dependencies using npm
	npm install
	touch $@

dev: .npminstall
	@echo Running webapp for development
	npm run dev

build: .npminstall
	@echo Building webapp
	npm run build

docker-build:
	@echo Building webapp with Docker
	docker run --rm -v $(PWD):/smsender -w /smsender node:6.9.5 sh -c "make build"

clean:
	@echo Cleaning webapp
	rm -f public/bundle.js
	rm -f .npminstall
	rm -rf node_modules