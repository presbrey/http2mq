package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/fcgi"
	"os"
	"runtime"

	"github.com/presbrey/http2mq/lib"
)

var (
	bind    = flag.String("bind", "", "bind address (empty=fcgi)")
	handler = http2mq.Handler{}
)

func init() {
	flag.Parse()
	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
}

func main() {
	var err error

	if bind == nil || len(*bind) == 0 {
		err = fcgi.Serve(nil, handler)
	} else {
		err = http.ListenAndServe(*bind, handler)
	}
	if err != nil {
		log.Fatalln(err)
	}
}
