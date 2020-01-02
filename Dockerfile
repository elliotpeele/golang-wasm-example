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

## Go build
ARG GO_VERSION=1.13.4
FROM golang:${GO_VERSION} AS gobuild
# cache deps
WORKDIR /gobuild
ENV GO111MODULE=on
ENV GOROOT=/usr/local/go
COPY go.mod go.sum ./
RUN go mod download
# build
COPY . .
RUN make

## Final output
FROM debian:stretch-slim
RUN groupadd -r wasm --gid=1000 && useradd -r -g wasm --uid=1000 wasm
COPY --from=gobuild /gobuild/build/linux/golang-wasm-example /usr/bin/golang-wasm-example
RUN apt-get update && apt-get install -y ca-certificates

USER wasm
CMD ["golang-wasm-example", "serve"]