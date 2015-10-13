package http2mq

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type confirmData struct {
	SubscribeURL string
}

func confirm(body []byte) bool {
	c := new(confirmData)

	if err := json.Unmarshal(body, c); err != nil {
		return false
	}

	if c.SubscribeURL == "" {
		return false
	}

	go func() {
		resp, err := http.Get(c.SubscribeURL)
		if err != nil {
			log.Println(err)
			return
		}
		b, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode > 299 {
			log.Println("confirm failed, code: %d, data: %s", resp.StatusCode, string(b))
			return
		}
		log.Println("confirm success, code: %d, data: %s", resp.StatusCode, string(b))
	}()

	return true
}
