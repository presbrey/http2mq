package http2mq

import (
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	XFF = "X-Forwarded-For"
)

var (
	autoConfirm   = flag.Bool("confirm", true, "follow SubscribeURLs")
	cors          = flag.Bool("cors", true, "allow CORS")
	escapeBody    = flag.Bool("escapeBody", false, "request body will be Go-escaped")
	escapeHeaders = flag.Bool("escapeHeaders", true, "request headers will be Go-escaped")
	setXFF        = flag.Bool("xForwardedFor", true, "prepend remote address to "+XFF)
	successCode   = flag.Int("successCode", 201, "onSuccess HTTP status code")
)

func escape(s string) (r string) {
	if len(s) > 0 {
		r = strconv.Quote(s)
		r = strings.Replace(r[1:len(r)-1], `\"`, `"`, -1)
	}
	return
}

type Handler struct{ http.Handler }

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var (
		data []byte
		err  error
		xff  string
	)

	if *cors {
		origins := req.Header["Origin"]
		if len(origins) > 0 {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Origin", origins[0])
		}
		if req.Method == "OPTIONS" {
			corsReqH := req.Header["Access-Control-Request-Headers"]
			if len(corsReqH) > 0 {
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(corsReqH, ", "))
			}
			corsReqM := req.Header["Access-Control-Request-Method"]
			if len(corsReqM) > 0 {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(corsReqM, ", "))
			} else {
				w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			}
			if len(origins) < 1 {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}
			w.WriteHeader(200)
			return
		}
	}

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

	if *autoConfirm && confirm(data) {
		w.WriteHeader(200)
		return
	}

	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}

	if *escapeBody {
		data = []byte(escape(string(data)))
	}

	head := amqp.Table{}
	if *escapeHeaders {
		for k, v := range req.Header {
			head[escape(k)] = escape(strings.Join(v, "\n"))
		}
	} else {
		for k, v := range req.Header {
			head[k] = strings.Join(v, "\n")
		}
	}

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

	if len(req.Host) > 0 {
		head["Host"] = req.Host
	}
	head["Time"] = fmt.Sprintf("%d", time.Now().Unix())
	elt := &Request{
		Headers: head,
		Body:    data,
	}
	incoming <- elt
	w.WriteHeader(*successCode)
}
