package main

import (
	"net/http"
)

func cors(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fs.ServeHTTP(w, r)
	}
}

func main() {
	DownloadWasmExec()
	fileServer := http.FileServer(http.Dir("."))
	http.Handle("/", cors(fileServer))
	println("Listening on port 3000...")
	http.ListenAndServe(":3000", nil)
}
