package main

import "net/http"

func main() {
    http.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
        w.Write([]byte("OK"))
    })
    http.ListenAndServe(":18800", nil)
}
