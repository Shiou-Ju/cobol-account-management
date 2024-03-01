package main

import (
	"fmt"
	"net/http"
)

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	http.HandleFunc("/", helloWorldHandler)
	fmt.Println("Server is running on port 3001")
	http.ListenAndServe(":3001", nil)
}
