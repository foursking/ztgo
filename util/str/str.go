package str

import (
	"bytes"
	"strconv"
	"strings"
	"sync"
)

var (
	bfPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer([]byte{})
		},
	}
)

// JoinInt64 format int64 slice like:n1,n2,n3.
func JoinInt64(is []int64) string {
	if len(is) == 0 {
		return ""
	}
	if len(is) == 1 {
		return strconv.FormatInt(is[0], 10)
	}
	buf := bfPool.Get().(*bytes.Buffer)
	for _, i := range is {
		buf.WriteString(strconv.FormatInt(i, 10))
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	s := buf.String()
	buf.Reset()
	bfPool.Put(buf)
	return s
}

// SplitInt64 split string into int64 slice.
func SplitInt64(s string) ([]int64, error) {
	if s == "" {
		return nil, nil
	}
	sArr := strings.Split(s, ",")
	res := make([]int64, 0, len(sArr))
	for _, sc := range sArr {
		i, err := strconv.ParseInt(sc, 10, 64)
		if err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}

// JoinInt format int slice like:n1,n2,n3.
func JoinInt(is []int) string {
	if len(is) == 0 {
		return ""
	}
	if len(is) == 1 {
		return strconv.Itoa(is[0])
	}
	buf := bfPool.Get().(*bytes.Buffer)
	for _, i := range is {
		buf.WriteString(strconv.Itoa(i))
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	s := buf.String()
	buf.Reset()
	bfPool.Put(buf)
	return s
}

// SplitInt split string into int slice.
func SplitInt(s string) ([]int, error) {
	if s == "" {
		return nil, nil
	}
	sArr := strings.Split(s, ",")
	res := make([]int, 0, len(sArr))
	for _, sc := range sArr {
		i, err := strconv.Atoi(sc)
		if err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}
