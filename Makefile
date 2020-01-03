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

TARGET      := golang-wasm-example
BUILD_ENV   := GO111MODULE=on
PACKAGE     := github.com/elliotpeele/golang-wasm-example
GO_OUT_PATH := build
TOPDIR      := $(PWD)
TOOLS       := $(TOPDIR)/tools

# build variables
DATE     = $(shell date +%Y-%m-%d)
VERSION  = $(shell git describe --tags --dirty 2>/dev/null || echo 0.0.0-0)
COMMIT   = $(shell git rev-parse --short HEAD 2>/dev/null)
LDFLAGS  = "-X $(PACKAGE)/cmd.commitHash=$(COMMIT) -X $(PACKAGE)/cmd.buildDate=$(DATE) -X $(PACKAGE)/cmd.version=$(VERSION)"

# commands
GO      = $(BUILD_ENV) go
GOBUILD = go build -v

.PHONY: all
binaries: frontend sampledata linux windows osx

# tools
GOBINDATA = $(TOOLS)/go-bindata
$(GOBINDATA):
	$(GO) build -o $@ github.com/shuLhan/go-bindata

.PHONY: tools
toos: | $(GOBINDATA)

.PHONY: clean
clean:
	rm -rf $(GO_OUT_PATH)
	rm -rf $(TOOLS)
	make -C frontend clean

	rm -f sampledata/gen
	rm -rf sampledata/generated
	rm -f sampledata/generated_data.go

.PHONY: sampledata
sampledata: tools
	PATH=$(TOOLS):$$PATH $(GO) generate -x ./sampledata

.PHONY: frontend
frontend:
	make -C frontend

.PHONY: linux
linux:
	mkdir -p $(GO_OUT_PATH)/linux
	$(BUILD_ENV) GOARCH=amd64 GOOS=linux $(GOBUILD) -ldflags $(LDFLAGS) -o $(GO_OUT_PATH)/linux/$(TARGET)

.PHONY: windows
windows:
	mkdir -p $(GO_OUT_PATH)/windows
	$(BUILD_ENV) GOARCH=amd64 GOOS=windows $(GOBUILD) -ldflags $(LDFLAGS) -o $(GO_OUT_PATH)/windows/$(TARGET)

.PHONY: osx
osx:
	mkdir -p $(GO_OUT_PATH)/osx
	$(BUILD_ENV) GOARCH=amd64 GOOS=darwin $(GOBUILD) -ldflags $(LDFLAGS) -o $(GO_OUT_PATH)/osx/$(TARGET)