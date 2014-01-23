package main

import (
	"github.com/drewolson/testflight"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGET(t *testing.T) {
	testflight.WithServer(handler, func(r *testflight.Requester) {
		response := r.Get("/?test=1")
		assert.Equal(t, 201, response.StatusCode)
	})
}

func TestPOST(t *testing.T) {
	testflight.WithServer(handler, func(r *testflight.Requester) {
		response := r.Post("/", testflight.JSON, `{"test":1}`)
		assert.Equal(t, 201, response.StatusCode)
	})
}

func TestPing(t *testing.T) {
	testflight.WithServer(handler, func(r *testflight.Requester) {
		response := r.Get("/ping")
		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "OK", response.Body)
	})
}
