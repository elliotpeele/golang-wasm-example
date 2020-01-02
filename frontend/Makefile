# Copyright (c) 2019 Elliot Peele <elliot@bentlogic.net>
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.

TARGET      := main.wasm
BUILD_ENV   := GO111MODULE=on
PACKAGE     := github.com/elliotpeele/golang-wasm-example/frontend
TOPDIR      := $(PWD)
TOOLS       := $(TOPDIR)/tools

# build variables
DATE     = $(shell date +%Y-%m-%d)
VERSION  = $(shell git describe --tags --dirty 2>/dev/null || echo 0.0.0-0)
COMMIT   = $(shell git rev-parse --short HEAD 2>/dev/null)
LDFLAGS  = "-X main.commitHash=$(COMMIT) -X main.buildDate=$(DATE) -X main.version=$(VERSION)"
PACKAGES = "./..."

# commands
GO      = $(BUILD_ENV) go
GOBUILD = go build -v

.PHONY: all
binaries: wasm wasm_exec generate

GOBINDATA = $(TOOLS)/go-bindata
$(GOBINDATA):
	$(GO) build -o $@ github.com/shuLhan/go-bindata/cmd/go-bindata
	chmod +x $@

.PHONY: clean
clean:
	rm -rf assets/dist/main.wasm
	rm -rf assets/dist/wasm_exec.js
	rm -rf assets/assets_generated.go

.PHONY: tools
tools: | $(GOBINDATA)

.PHONY: generate
generate: tools 
	PATH=$(TOOLS):$$PATH $(GO) generate -x ./assets

.PHONY: wasm
wasm:
	$(BUILD_ENV) GOARCH=wasm GOOS=js $(GOBUILD) -ldflags $(LDFLAGS) -o assets/dist/$(TARGET)

.PHONY: wasm_exec
wasm_exec:
	cp $(GOROOT)/misc/wasm/wasm_exec.js assets/dist/