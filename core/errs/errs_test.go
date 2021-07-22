package errs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErr_WithDetail(t *testing.T) {
	e1 := New(8, "test WithDetail").WithDetail("msg")
	assert.Equal(t, map[string]string{"detail": "msg"}, e1.Details)
	e2 := New(8, "test WithDetail").WithDetail("key", "msg")
	assert.Equal(t, map[string]string{"key": "msg"}, e2.Details)
	e3 := New(8, "test WithDetail").WithDetail("key1", "msg1", "key2", "msg2")
	assert.Equal(t, map[string]string{"key1": "msg1"}, e3.Details)
}
