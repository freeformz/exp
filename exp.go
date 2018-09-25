package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"
)

var version string

func debugForm(r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println("Error Parsing Form: ", err)
		return
	}
	fmt.Println("Form Values:")
	if len(r.Form) == 0 {
		fmt.Printf("No Form Values")
	}
	for k := range r.Form {
		fmt.Printf("\t%q: %q\n", k, r.FormValue(k))
	}
}

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
		fmt.Printf("version=%q go-version=%q numcpu=%q\n", version, runtime.Version(), runtime.NumCPU())
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rdrops := r.Header.Get("Logshuttle-Drops")
		d, err := strconv.Atoi(rdrops)
		if err == nil {
			atomic.AddUint64(&drops, uint64(d))
		}

		if debug {
			defer func() {
				fmt.Println("---")
			}()
			fmt.Println(time.Now())
			fmt.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)
			fmt.Printf("Scheme: %s\n", r.URL.Scheme)
			fmt.Printf("Host: %s\n", r.Host)
			fmt.Println("Headers:")
			for k := range r.Header {
				fmt.Printf("\t%s: %s\n", k, r.Header.Get(k))
			}
			debugForm(r)
			fmt.Printf("Content Length: %d\n", r.ContentLength)
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error Reading Body:", err)
			return
		}
		if debug {
			fmt.Println("Body:")
			fmt.Println(string(body))
		}
	})

	http.ListenAndServe(listen, nil)
}
