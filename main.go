package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
	})

	fmt.Println("Server is running on port 3000")
	http.ListenAndServe(":3000", nil)
}
