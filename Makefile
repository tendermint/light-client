DOC_PKGS:=./certifiers ./extensions ./extensions/basecoin ./proxy ./proxy/types ./tx
REPO:=github.com/tendermint/light-client

.PHONY: install build test list_pkg docs clean_docs get_vendor_deps tools $(DOC_PKGS)

install: get_vendor_deps
	go install ./cmd/...

build:
	go build ./cmd/...

test/shunit2:
	wget "https://raw.githubusercontent.com/kward/shunit2/master/source/2.1/src/shunit2" \
		-q -O test/shunit2

test_cli: test/shunit2
	@./test/keys.sh
	@./test/init.sh

test: build test_unit test_cli

# note that we start tendermint nodes in rpc/tests and extensions/basecoin
# we cannot currently run these tests in parallel
test_unit:
	go test -p 1 `glide novendor`

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
	#@go get github.com/davecheney/godoc2md
	@go get github.com/Masterminds/glide
