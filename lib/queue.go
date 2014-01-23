package http2mq

import (
	"flag"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type Request struct {
	Headers     amqp.Table
	ContentType string
	Body        []byte
}
type RequestChan chan *Request

var (
	backlog  = flag.Uint("backlog", 8192, "incoming channel capacity")
	backoff  = flag.Duration("backoff", time.Second, "pause between errors")
	exchange = flag.String("exchange", "test", "AMQP exchange name")
	key      = flag.String("key", "test", "AMQP routing key")
	tag      = flag.String("tag", "http2mq", "AMQP consumer tag")
	uri      = flag.String("uri", "amqp://localhost:5672/", "AMQP URI")

	incoming = make(RequestChan, *backlog)
)

func dial() (ch *amqp.Channel, err error) {
	conn, err := amqp.Dial(*uri)
	if err != nil {
		return
	}

	channel, err := conn.Channel()
	if err != nil {
		return
	}

	err = channel.ExchangeDeclare(
		*exchange, // name
		"direct",  // type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return
	}

	ch = channel
	return
}

func (rc RequestChan) publish() error {
	var (
		err error
		out *amqp.Channel
	)
	for r := range rc {
		if r == nil {
			break
		}
		for out == nil {
			out, err = dial()
			if err != nil {
				out = nil
				log.Println(err)
				time.Sleep(*backoff)
			}
		}

		err = out.Publish(*exchange, *key, false, false,
			amqp.Publishing{
				Headers:      amqp.Table(r.Headers),
				ContentType:  r.ContentType,
				DeliveryMode: amqp.Persistent,
				Body:         r.Body,
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	go func() {
		for {
			err := incoming.publish()
			if err != nil {
				log.Println(err)
			}
			time.Sleep(*backoff)
		}
	}()
}
