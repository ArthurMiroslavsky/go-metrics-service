package main

import (
	"fmt"
	"net/http"
)

func updateHandler(w http.ResponseWriter, r *http.Request) {
	u := r.URL
	fmt.Println(u)
}

func main() {
	http.HandleFunc("/update/", updateHandler)
	http.ListenAndServe("127.0.0.1:8080", nil)
}
