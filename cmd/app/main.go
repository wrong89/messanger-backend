package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from test handler"))
	})

	http.ListenAndServe(":9090", nil)
	fmt.Println("Hello World")
}
