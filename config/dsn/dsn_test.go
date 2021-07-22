package dsn

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	dsnstr := "udp://username:password@192.168.3.5:8000,192.168.3.5:9000?timeout=1s&foo.bar=1"
	d, err := Parse(dsnstr)
	assert.Nil(t, err)
	assert.Equal(t, "udp", d.Scheme)
	assert.Equal(t, "username", d.User.Username())
	pwd, isSet := d.User.Password()
	assert.True(t, isSet)
	assert.Equal(t, "password", pwd)
	assert.Equal(t, "192.168.3.5:8000,192.168.3.5:9000", d.Host)
	q := d.Query()
	assert.Equal(t, "1s", q.Get("timeout"))
	assert.Equal(t, "1", q.Get("foo.bar"))
}

func TestDSN_Bind(t *testing.T) {
	type Foo struct {
		Bar int `dsn:"query.bar"`
	}
	type Config struct {
		Scheme    string        `dsn:"scheme" validate:"required"`
		Addresses []string      `dsn:"address"`
		Username  string        `dsn:"username"`
		Password  string        `dsn:"password"`
		Timeout   time.Duration `dsn:"query.timeout"`
		Foo       Foo           `dsn:"query.foo"`
		Def       string        `dsn:"query.def,hello"`
	}
	dsnstr := "udp://username:password@192.168.3.5:8000,192.168.3.5:9000?timeout=10s&foo.bar=1"
	d, err := Parse(dsnstr)
	assert.Nil(t, err)
	c := &Config{}
	err = d.Bind(c)
	assert.Nil(t, err)
	fmt.Printf("%+v error(%v)", c, err)
}
