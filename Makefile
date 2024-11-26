VERSION := 0.1
REV := $(shell git rev-parse --short HEAD)
SRC := $(filter-out $(wildcard *_test.go), $(wildcard *.go))

GOBUILDFLAGS=-v -trimpath
LDFLAGS=-w -s -X main.Version=$(VERSION).$(REV)
CGO := 0

.PHONY: test
test: paintface.prg
	open $<

retrospex: $(SRC)
	go build $(GOBUILDFLAGS) -ldflags="$(LDFLAGS)" -o "$@"

.PHONY: all
all: \
	retrospex_macos_arm64.zip \
	retrospex_macos_amd64.zip \
	retrospex_linux_arm64.zip \
	retrospex_linux_amd64.zip \
	retrospex_windows_arm64.exe.zip \
	retrospex_windows_amd64.exe.zip \
	retrospex_windows_x86.exe.zip

.PHONY: clean
clean:
	rm *.png || true
	rm *.zip || true
	rm *.prg || true
	rm retrospex_*

.PHONY: install
install: $(SRC)
	go build -o "${HOME}/bin/retrospex"

%.zip: %
	zip -m -9 $@ $<

retrospex_linux_amd64: $(SRC)
	CGO_ENABLED=$(CGO) GOOS=linux GOARCH=amd64 go build $(GOBUILDFLAGS) -ldflags="$(LDFLAGS) -X main.Arch=linux.amd64" -o $@

retrospex_linux_arm64: $(SRC)
	CGO_ENABLED=$(CGO) GOOS=linux GOARCH=arm64 go build $(GOBUILDFLAGS) -ldflags="$(LDFLAGS) -X main.Arch=linux.arm64" -o $@

retrospex_macos_amd64: $(SRC)
	CGO_ENABLED=$(CGO) GOOS=darwin GOARCH=amd64 go build $(GOBUILDFLAGS) -ldflags="$(LDFLAGS) -X main.Arch=macos.amd64" -o $@

retrospex_macos_arm64: $(SRC)
	CGO_ENABLED=$(CGO) GOOS=darwin GOARCH=arm64 go build $(GOBUILDFLAGS) -ldflags="$(LDFLAGS) -X main.Arch=macos.arm64" -o $@

retrospex_windows_amd64.exe: $(SRC)
	CGO_ENABLED=$(CGO) GOOS=windows GOARCH=amd64 go build $(GOBUILDFLAGS) -ldflags="$(LDFLAGS) -X main.Arch=windows.amd64" -o $@

retrospex_windows_arm64.exe: $(SRC)
	CGO_ENABLED=$(CGO) GOOS=windows GOARCH=arm64 go build $(GOBUILDFLAGS) -ldflags="$(LDFLAGS) -X main.Arch=windows.arm64" -o $@

retrospex_windows_x86.exe: $(SRC) 
	CGO_ENABLED=$(CGO) GOOS=windows GOARCH=386 go build $(GOBUILDFLAGS) -ldflags="$(LDFLAGS) -X main.Arch=windows.x86" -o $@

include scripts/scsprites.mk
