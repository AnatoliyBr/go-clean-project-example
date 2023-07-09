.PHONY: build
build: 
		go build -v ./cmd/app

.PHONY: test
test:
		go test -v -race ./...

.PHONY: upload
upload:
		git add -A
		git commit -m "$m"
		git push origin

.PHONY: download
download:
		git pull

.PHONY: tag
tag:	
		git tag -a $v -m "$m"
		git push origin --tags

.DEFAULT_GOAL := build