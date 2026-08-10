package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jackob2004/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	ping(res, req)

	assert.Equal(t, res.Code, http.StatusOK)
	assert.Equal(t, res.Body.String(), "OK")
}
