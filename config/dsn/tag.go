package dsn

import "strings"

const tagName = "dsn"

type tag struct {
	Name    string
	Default string
}

func newTag(s string) *tag {
	ss := strings.SplitN(s, ",", 2)
	t := tag{Name: ss[0]}
	if len(ss) == 2 {
		t.Default = ss[1]
	}
	return &t
}
