package main

import (
	"net/http"
)

func main() {

	server := http.NewServeMux()

	server.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {})

	if err := http.ListenAndServe(":3000", server); err != nil {
		panic(err)
	}
}
