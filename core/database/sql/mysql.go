package sql

import (
	"github.com/foursking/ztgo/log"

	_ "github.com/go-sql-driver/mysql"
)

// NewMySQL new db and retry connection when has error.
func NewMySQL(opts ...Option) (db *DB) {
	options := applyOptions(opts...)
	db, err := Open(options)
	if err != nil {
		log.Fatalf("open mysql error(%v)", err)
	}
	return
}
