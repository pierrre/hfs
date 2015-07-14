// Package hfs provides a HTTP File Server.
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

var (
	flagHTTP = ":8080"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	parseFlags()
	p := getPath()
	h, err := getHTTPHandler(p)
	if err != nil {
		return err
	}
	fmt.Printf("serving %s on %s\n", p, flagHTTP)
	return http.ListenAndServe(flagHTTP, h)
}

func parseFlags() {
	flag.StringVar(&flagHTTP, "http", flagHTTP, "HTTP addr")
	flag.Parse()
}

func getPath() string {
	p := flag.Arg(0)
	if len(p) == 0 {
		p = "."
	}
	return p
}

func getHTTPHandler(p string) (http.Handler, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	var h http.Handler
	if fi.IsDir() {
		h = http.FileServer(http.Dir(p))
	} else {
		h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, p)
		})
	}
	return h, nil
}
