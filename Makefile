PKGS := github.com/pkg/errors
SRCDIRS := $(shell go list -f '{{.Dir}}' $(PKGS))
GO := go

check: test vet gofmt misspell unconvert staticcheck ineffassign unparam

test: 
	$(GO) test $(PKGS)

vet: | test
	$(GO) vet $(PKGS)

staticcheck:
	$(GO) install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck -checks all $(PKGS)

misspell:
	$(GO) install github.com/client9/misspell/cmd/misspell@latest
	misspell \
		-locale GB \
		-error \
		*.md *.go

unconvert:
	$(GO) install github.com/mdempsky/unconvert@latest
	unconvert -v $(PKGS)

ineffassign:
	$(GO) install github.com/gordonklaus/ineffassign@latest
	find $(SRCDIRS) -name '*.go' | xargs ineffassign

pedantic: check errcheck

unparam:
	$(GO) install mvdan.cc/unparam@latest
	unparam ./...

errcheck:
	$(GO) install github.com/kisielk/errcheck@latest
	errcheck $(PKGS)

gofmt:  
	@echo Checking code is gofmted
	@test -z "$(shell gofmt -s -l -d -e $(SRCDIRS) | tee /dev/stderr)"
