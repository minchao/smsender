.PHONY: build dev clean

.npminstall: package.json
	@echo Getting dependencies using npm
	npm install
	touch $@

build: .npminstall
	@echo Building webapp
	npm run build

dev: .npminstall
	@echo Running webapp for development
	npm run dev

clean:
	@echo Cleaning webapp
	rm -f public/bundle.js
	rm -f .npminstall
	rm -rf node_modules