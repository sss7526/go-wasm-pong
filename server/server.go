package main

import (
    "log"
    "net/http"
)

func main() {
    // Serve the HTML file and Wasm content
    fs := http.FileServer(http.Dir("../webapp"))
    http.Handle("/", fs)

    log.Println("Serving on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}