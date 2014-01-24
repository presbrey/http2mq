package http2mq

import (
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	setXFF = flag.Bool("xForwardedFor", true, "prepend remote address to X-Forwaded-For")
)

const (
	XFF = "X-Forwarded-For"
)

type Handler struct{ http.Handler }

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var (
		data []byte
		err  error
		xff  string
	)

	if len(req.URL.Path) >= 5 && req.URL.Path == "/ping" {
		fmt.Fprint(w, "OK")
		return
	}

	if req.Method == "GET" {
		data = []byte(req.URL.RawQuery)
	} else {
		if data, err = ioutil.ReadAll(req.Body); err != nil {
			log.Println(err)
		}
	}

	if len(data) == 0 {
		w.WriteHeader(204)
		return
	}

	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}

	head := amqp.Table{}
	for k, v := range req.Header {
		head[k] = strings.Join(v, "\n")
	}
	ctype := req.Header.Get("Content-Type")

	if *setXFF {
		xff_ := req.Header.Get(XFF)
		addr, err := net.ResolveTCPAddr("tcp", req.RemoteAddr)
		if err == nil {
			xff = addr.IP.String()
		} else {
			xff = req.RemoteAddr
		}
		if len(xff) > 0 && len(xff_) > 0 {
			xff += ", "
		}
		xff += xff_
		if len(xff) > 0 {
			head[XFF] = xff
		}
	}

	head["Time"] = fmt.Sprintf("%d", time.Now().Unix())
	elt := &Request{
		Headers:     head,
		ContentType: ctype,
		Body:        data,
	}
	incoming <- elt
	w.WriteHeader(201)
}
