GoLang WASM Example
===================

This is an example project showing [Go WASM](https://github.com/golang/go/wiki/WebAssembly) in action. (aka How a backend developer that knows very little about javascript can write a full frontend interface in a few days)

I am primarily a backend service developer, I needed a frontend for a project I was working on and ran across WASM support in Go. After much googling and reading, I managed to get a usable interface for displaying tables of data with searching and pagination with a few days of work.

Thanks to the [wasmws project](https://github.com/tarndt/wasmws) for providing a way to run [grpc](https://grpc.io/) over [websockets](https://en.wikipedia.org/wiki/WebSocket).

### Thing to know before going on the WASM adventure:
* WASM support in Go is still relatively new and should be considered experimental
* Support was originally introduced in Go 1.11
* The syscall/js API changed between Go 1.12 and 1.13 in a way that broke a lot of the examples and libraries that are around
* I haven't been able to get the WASM to load on mobile browsers
* There is a delay in the WASM load, at least in chrome. There is probably a way to improve this, but I haven't spent time on it.

### Walkthrough
There is a backend service that provides a websocket wrapped grpc service as well as static file serving. Take a look at cmd/serve.go for the setup and api/api.go for the API.

There are two parts to the frontend, a WASM binary and an index.html that loads the WASM and stores elements for use. The frontend is a grpc/websocket client that registers a set of methods into the dom. When buttons are clicked, the methods are called, rendering the data that is retrieved from the backend service.

The index.html has a hidden div that contains several template elements that are used for the data display, pagination, and searching.

### Building
Run `docker build . -t wasm` or `make`

### Running
#### Local build
```sh
./build/linux/golang-wasm-example serve --port 8080
```
#### Docker
```sh
docker run -p 8080:8080 wasm
```