# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: gvbc android ios gvbc-cross evm all test clean
.PHONY: gvbc-linux gvbc-linux-386 gvbc-linux-amd64 gvbc-linux-mips64 gvbc-linux-mips64le
.PHONY: gvbc-linux-arm gvbc-linux-arm-5 gvbc-linux-arm-6 gvbc-linux-arm-7 gvbc-linux-arm64
.PHONY: gvbc-darwin gvbc-darwin-386 gvbc-darwin-amd64
.PHONY: gvbc-windows gvbc-windows-386 gvbc-windows-amd64

GOBIN = ./build/bin
GO ?= latest
GORUN = env GO111MODULE=on go run

gvbc:
	$(GORUN) build/ci.go install ./cmd/gvbc
	@echo "Done building."
	@echo "Run \"$(GOBIN)/gvbc\" to launch gvbc."

all:
	$(GORUN) build/ci.go install

android:
	$(GORUN) build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/gvbc.aar\" to use the library."
	@echo "Import \"$(GOBIN)/gvbc-sources.jar\" to add javadocs"
	@echo "For more info see https://stackoverflow.com/questions/20994336/android-studio-how-to-attach-javadoc"
	
ios:
	$(GORUN) build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/Gvbc.framework\" to use the library."

test: all
	$(GORUN) build/ci.go test

lint: ## Run linters.
	$(GORUN) build/ci.go lint

clean:
	env GO111MODULE=on go clean -cache
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go install ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

# Cross Compilation Targets (xgo)

gvbc-cross: gvbc-linux gvbc-darwin gvbc-windows gvbc-android gvbc-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/gvbc-*

gvbc-linux: gvbc-linux-386 gvbc-linux-amd64 gvbc-linux-arm gvbc-linux-mips64 gvbc-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/gvbc-linux-*

gvbc-linux-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/gvbc
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/gvbc-linux-* | grep 386

gvbc-linux-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/gvbc
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gvbc-linux-* | grep amd64

gvbc-linux-arm: gvbc-linux-arm-5 gvbc-linux-arm-6 gvbc-linux-arm-7 gvbc-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/gvbc-linux-* | grep arm

gvbc-linux-arm-5:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/gvbc
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/gvbc-linux-* | grep arm-5

gvbc-linux-arm-6:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/gvbc
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/gvbc-linux-* | grep arm-6

gvbc-linux-arm-7:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/gvbc
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/gvbc-linux-* | grep arm-7

gvbc-linux-arm64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/gvbc
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/gvbc-linux-* | grep arm64

gvbc-linux-mips:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/gvbc
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/gvbc-linux-* | grep mips

gvbc-linux-mipsle:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/gvbc
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/gvbc-linux-* | grep mipsle

gvbc-linux-mips64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/gvbc
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/gvbc-linux-* | grep mips64

gvbc-linux-mips64le:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/gvbc
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/gvbc-linux-* | grep mips64le

gvbc-darwin: gvbc-darwin-386 gvbc-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/gvbc-darwin-*

gvbc-darwin-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/gvbc
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/gvbc-darwin-* | grep 386

gvbc-darwin-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/gvbc
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gvbc-darwin-* | grep amd64

gvbc-windows: gvbc-windows-386 gvbc-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/gvbc-windows-*

gvbc-windows-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/gvbc
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/gvbc-windows-* | grep 386

gvbc-windows-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/gvbc
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gvbc-windows-* | grep amd64

goreleaser-install:
	go install github.com/goreleaser/goreleaser@latest

goreleaser-release:
	goreleaser release --clean

goreleaser-build:
	goreleaser build --clean --snapshot
