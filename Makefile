install: get_vendor_deps
	go install ./cmd/...

test:
	go test `glide novendor`

get_vendor_deps:
	go get github.com/Masterminds/glide
	rm -rf vendor/
	glide install
