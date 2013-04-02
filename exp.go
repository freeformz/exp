package main

import (
	"fmt"
	"time"
	"net/http"
	"os"
)

func main() {
	listen := os.Getenv("LISTEN")

	if listen == "" {
		listen = ":8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(time.Now())
		fmt.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)
		fmt.Println("Headers:")
		for k, _ := range r.Header {
			fmt.Printf("\t%s: %s\n", k, r.Header.Get(k))
		}
		fmt.Printf("Content Length: %d\n", r.ContentLength)
		if r.ContentLength > 0 {
			var body = make([]byte, r.ContentLength)
			r.Body.Read(body)
			fmt.Println("Body:")
			fmt.Println(string(body))
		}
		w.WriteHeader(http.StatusOK)
		fmt.Println()
	})

	http.ListenAndServe(listen, nil)
}
