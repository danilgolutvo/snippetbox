package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {

	app := newTestApp(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}
