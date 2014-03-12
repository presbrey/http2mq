# http2mq [![Build Status](https://travis-ci.org/presbrey/http2mq.png)](https://travis-ci.org/presbrey/http2mq)

## Install

Use the `go get` command eg.

    go get github.com/presbrey/http2mq

## Usage

`http2mq` accepts command-line arguments:
~~~
  -backlog=8192: incoming channel capacity
  -backoff=1s: pause between errors
  -bind="": bind address (empty=fcgi)
  -exchange="test": AMQP exchange name
  -exchangeType="fanout": AMQP exchange type
  -key="test": AMQP routing key
  -successCode=201: onSuccess HTTP status code
  -tag="http2mq": AMQP consumer tag
  -uri="amqp://localhost:5672/": AMQP URI
  -xForwardedFor=true: prepend remote address to X-Forwaded-For
~~~

## License

[MIT](http://joe.mit-license.org/)
