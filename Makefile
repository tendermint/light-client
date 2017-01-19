
install: get_vendor_deps
	go install ./cmd/...

build:
	go build ./cmd/...

test: build
	go test `glide novendor`

get_vendor_deps:
	go get github.com/Masterminds/glide
	rm -rf vendor/
	glide install
