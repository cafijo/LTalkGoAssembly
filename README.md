# Lightning talks Go Assembly Example Projects
Some small examples of using Go and WebAssembly
# Usage
The `server` directory contains a very small HTTP server implementation that 
simply hosts the files of the current working directory.
It also downloads the `wasm_exec.js` JavaScript bridge from the official golang
repository, if the file doesn't already exist.

```bash
export APIKEY_CHATGPT='YOUR_API_KEY_CHATGPT'
go build -o server.bin ./server
./server.bin
```
You can also use any other web-server that provides static file hosting.

## Compiling:
Browse to http://localhost:3000 after building the example:

```bash
GOARCH=wasm GOOS=js go build -o test.wasm ./chatgpt
```

