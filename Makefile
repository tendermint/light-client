.PHONY: install build test get_vendor_deps tools

build: tools
	go build `glide novendor`

test:
	go test -p 1 `glide novendor`

get_vendor_deps: tools
	@rm -rf vendor/
	glide install

tools:
	@go get github.com/Masterminds/glide
