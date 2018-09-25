package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var version string

var mu sync.Mutex
var n int

func p(s string) {
	mu.Lock()
	defer func() {
		n++
		mu.Unlock()
	}()
	if s[0] == '\t' {
		fmt.Printf("%d: \t%q\n", n, s[1:])
		return
	}
	fmt.Printf("%d: %q\n", n, s)
}

func debugForm(r *http.Request) {
	if err := r.ParseForm(); err != nil {
		p("Error Parsing Form: " + err.Error())
		return
	}
	p("Form Values:")
	if len(r.Form) == 0 {
		p("\tNo Form Values")
	}
	for k := range r.Form {
		p(fmt.Sprintf("\t%s: %s", k, r.FormValue(k)))
	}
}

func main() {
	listen := os.Getenv("LISTEN")
	debug := os.Getenv("DEBUG") != ""

	var drops uint64

	if listen == "" {
		listen = ":8080"
	}

	p("Listening on: " + listen)

	if debug {
		for _, e := range os.Environ() {
			p(e)
		}
		p(fmt.Sprintf("version=%s go-version=%s numcpu=%d", version, runtime.Version(), runtime.NumCPU()))
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rdrops := r.Header.Get("Logshuttle-Drops")
		d, err := strconv.Atoi(rdrops)
		if err == nil {
			atomic.AddUint64(&drops, uint64(d))
		}

		if debug {
			defer func() {
				p("---")
			}()
			p(time.Now().String())
			p(fmt.Sprintf("%s %s %s", r.Method, r.URL, r.Proto))
			p(fmt.Sprintf("Scheme: %s", r.URL.Scheme))
			p(fmt.Sprintf("Host: %s", r.Host))
			p(fmt.Sprintf("Headers:"))
			for k := range r.Header {
				p(fmt.Sprintf("\t%s: %s", k, r.Header.Get(k)))
			}
			debugForm(r)
			p(fmt.Sprintf("Content Length: %d\n", r.ContentLength))
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			p("Error Reading Body: " + err.Error())
			return
		}
		if debug {
			p("Body:")
			p(string(body))
		}
	})

	http.ListenAndServe(listen, nil)
}
