package http

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestServer_Recovery(t *testing.T) {
	s := NewServer()
	s.GET("/test/recovery", testRecovery)
	s.Run()
	time.Sleep(10 * time.Millisecond)
	rsp, err := http.Get("http://127.0.0.1:8080/test/recovery")
	assert.Nil(t, err)
	assert.Equal(t, 500, rsp.StatusCode)
}

func testRecovery(ctx *gin.Context) {
	panic(errors.New("test panic recovery"))
}
