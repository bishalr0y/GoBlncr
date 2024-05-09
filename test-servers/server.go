package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world")
	})
	fmt.Println("sever running at port 3000")
	http.ListenAndServe(":3000", nil)
}
