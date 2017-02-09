DOC_PKGS:=./cryptostore ./mock ./proxy ./rpc ./rpc/tests ./storage ./storage/filestorage ./storage/memstorage ./tx ./util
REPO:=github.com/tendermint/light-client

.PHONY: install build test list_pkg docs clean_docs get_vendor_deps tools $(DOC_PKGS)

install: get_vendor_deps
	go install ./cmd/...

build:
	go build ./cmd/...

test: build
	go test `glide novendor`

# run list_pkg manually to make DOC_PKGS -> Makefile won that fight
list_pkg:
	@find . -maxdepth 2 -type d ! -path './vendor*' ! -path './.*' ! -path './docs*' ! -path './cmd*' -exec echo {} \; | tr '\n' ' '; echo

# separated this out from the rest of the packages for naming issues
lightclient:
	godoc2md $(REPO) > docs/lightclient.md

$(DOC_PKGS):
	godoc2md $(REPO)/$@ > docs/$(notdir $@).md

docs: tools clean_docs lightclient $(DOC_PKGS)

clean_docs:
	@rm -rf docs
	@mkdir docs
# @echo $(DOC_PKGS)
#	for dir in `${DOC_PKGS}`; do godoc2md ${dir} > docs/`basename ${dir}`.md; done

get_vendor_deps: tools
	@rm -rf vendor/
	glide install

tools:
	@go get github.com/davecheney/godoc2md
	@go get github.com/Masterminds/glide
