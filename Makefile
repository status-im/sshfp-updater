VERSION =    $(shell cat VERSION)
BUILD_DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_HASH =   $(shell git rev-parse HEAD)

PACKAGE = sshfp-updater
BUILDS = linux-amd64 linux-386 linux-arm64 linux-mips64 windows-amd64.exe freebsd-amd64 darwin-amd64 darwin-arm64
BINARIES = $(addprefix bin/$(PACKAGE)-$(VERSION)-, $(BUILDS))
CHECKSUMS = bin/$(PACKAGE)-$(VERSION).sha256

LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildDate=$(BUILD_DATE) -X main.GitHash=$(GIT_HASH)"

$(PACKAGE): cmd/$(PACKAGE)
	go build -v -o $@ ./$^

release: $(BINARIES) checksums
checksums: $(CHECKSUMS)

bin:
	mkdir $@

bin/$(PACKAGE)-$(VERSION)-linux-%: bin
	env GOOS=linux GOARCH=$* CGO_ENABLED=0 go build $(LDFLAGS) -o $@ ./cmd/$(PACKAGE)

bin/$(PACKAGE)-$(VERSION)-darwin-%: bin
	env GOOS=darwin GOARCH=$* CGO_ENABLED=0 go build $(LDFLAGS) -o $@ ./cmd/$(PACKAGE)

bin/$(PACKAGE)-$(VERSION)-windows-%.exe: bin
	env GOOS=windows GOARCH=$* CGO_ENABLED=0 go build $(LDFLAGS) -o $@ ./cmd/$(PACKAGE)

bin/$(PACKAGE)-$(VERSION)-freebsd-%: bin
	env GOOS=freebsd GOARCH=$* CGO_ENABLED=0 go build $(LDFLAGS) -o $@ ./cmd/$(PACKAGE)

$(CHECKSUMS):
	sha256sum $(BINARIES) > $@

compile-analysis: cmd/$(PACKAGE)
	go build -gcflags '-m' ./$^

code-quality:
	-go vet ./cmd/$(PACKAGE)
	-gofmt -s -d ./cmd/$(PACKAGE)
	-golint ./cmd/$(PACKAGE)
	-gocyclo ./cmd/$(PACKAGE)
	-ineffassign ./cmd/$(PACKAGE)

test:
	go test -v ./cmd/$(PACKAGE)
