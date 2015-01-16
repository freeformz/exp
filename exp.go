package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

func main() {
	listen := os.Getenv("LISTEN")
	debug := os.Getenv("DEBUG") != ""

	var drops uint64

	if listen == "" {
		listen = ":8080"
	}

	fmt.Println("Listening on: " + listen)

	if debug {
		for _, e := range os.Environ() {
			fmt.Println(e)
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if debug {
			fmt.Println(time.Now())
			fmt.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)
			fmt.Println("Headers:")
		}
		rdrops := r.Header.Get("Logshuttle-Drops")
		d, err := strconv.Atoi(rdrops)
		if err == nil {
			atomic.AddUint64(&drops, uint64(d))
		}
		if debug {
			for k, _ := range r.Header {
				fmt.Printf("\t%s: %s\n", k, r.Header.Get(k))
			}
		}
		if debug {
			fmt.Printf("Content Length: %d\n", r.ContentLength)
		}

		if r.ContentLength > 0 {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Println("Error Reading Body:", err)
			} else {
				if debug {
					fmt.Println("Body:")
					fmt.Println(string(body))
				}
			}
		}
		if debug {
			fmt.Println()
		}
	})

	http.ListenAndServe(listen, nil)
}
