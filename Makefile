GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

build: fmtcheck
	go install ./cmd/statefaker

fmtcheck:
	@sh -c "'$(CURDIR)/gofmtcheck.sh'"

fmt:
	gofmt -w $(GOFMT_FILES)

test: fmtcheck
	go test -v ./...
