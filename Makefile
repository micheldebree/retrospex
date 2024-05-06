VERSION := 0.0
REV := $(shell git rev-parse --short HEAD)
SRC := $(filter-out $(wildcard *_test.go), $(wildcard *.go))

GOBUILDFLAGS=-v -trimpath
LDFLAGS=-w -s -X main.Version=$(VERSION).$(REV)
CGO := 0

# IMG := paintface.jpg
# IMG := colors2.jpg
# IMG := whitney-norm.png
# IMG := madonna01.png
IMG := whitney.png

.PHONY: test
test: $(SRC)
	go run $^ -o $@-$(REV).png $(IMG) && open $@-$(REV).png
	png2prg -v -d $@-$(REV).png
	open $@-$(REV).prg

.PHONY: all
all: \
	retrospex_macos_arm64 \
	retrospex_macos_amd64 \
	retrospex_linux_arm64 \
	retrospex_linux_amd64 \
	retrospex_windows_arm64.exe \
	retrospex_windows_amd64.exe \
	retrospex_windows_x86.exe

.PHONY: clean
clean:
	rm *.png || true
	rm retrospex_*

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

