package sql

import (
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/foursking/ztgo/errs"
	"github.com/foursking/ztgo/log"
	"github.com/foursking/ztgo/net/util/breaker"

	"github.com/go-sql-driver/mysql"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pkg/errors"
)

const (
	_family          = "sql_client"
	_slowLogDuration = time.Millisecond * 250
)

var (
	// ErrStmtNil prepared stmt error
	ErrStmtNil = errors.New("sql: prepare failed and stmt nil")
	// ErrNoMaster is returned by Master when call master multiple times.
	ErrNoMaster = errors.New("sql: no master instance")
	// ErrNoRows is returned by Scan when QueryRow doesn't return a row.
	// In such a case, QueryRow returns a placeholder *Row value that defers
	// this error until a Scan.
	ErrNoRows = sql.ErrNoRows
	// ErrTxDone transaction done.
	ErrTxDone = sql.ErrTxDone
)

// DB database.
type DB struct {
	write  *conn
	read   []*conn
	idx    int64
	master *DB
}

// conn database connection
type conn struct {
	*sql.DB
	breaker breaker.Breaker
	opt     *Options
	addr    string
}

// Tx transaction.
type Tx struct {
	db     *conn
	tx     *sql.Tx
	span   opentracing.Span
	ctx    context.Context
	cancel func()
}

// Row row.
type Row struct {
	err error
	*sql.Row
	db     *conn
	query  string
	args   []interface{}
	span   opentracing.Span
	cancel func()
}

// Scan copies the columns from the matched row into the values pointed at by dest.
func (r *Row) Scan(dest ...interface{}) (err error) {
	defer slowLog(fmt.Sprintf("Scan query(%s) args(%+v)", r.query, r.args), time.Now())
	if r.span != nil {
		defer r.span.Finish()
	}
	if r.err != nil {
		err = r.err
	} else if r.Row == nil {
		err = ErrStmtNil
	}
	if err != nil {
		return
	}
	err = r.Row.Scan(dest...)
	if r.cancel != nil {
		r.cancel()
	}
	r.db.onBreaker(&err)
	if err != ErrNoRows {
		err = errors.Wrapf(err, "query %s args %+v", r.query, r.args)
	}
	return
}

// Rows rows.
type Rows struct {
	*sql.Rows
	cancel func()
}

// Close closes the Rows, preventing further enumeration. If Next is called
// and returns false and there are no further result sets,
// the Rows are closed automatically and it will suffice to check the
// result of Err. Close is idempotent and does not affect the result of Err.
func (rs *Rows) Close() (err error) {
	err = errors.WithStack(rs.Rows.Close())
	if rs.cancel != nil {
		rs.cancel()
	}
	return
}

// Stmt prepared stmt.
type Stmt struct {
	db    *conn
	tx    bool
	query string
	stmt  atomic.Value
	span  opentracing.Span
}

// Open opens a database specified by its database driver name and a
// driver-specific data source name, usually consisting of at least a database
// name and connection information.
func Open(opt *Options) (*DB, error) {
	db := new(DB)
	d, err := connect(opt, opt.DSN)
	if err != nil {
		return nil, err
	}
	addr := parseDSNAddr(opt.DSN)
	brkGroup := breaker.NewGroup(opt.Breaker)
	brk := brkGroup.Get(addr)
	w := &conn{DB: d, breaker: brk, opt: opt, addr: addr}
	rs := make([]*conn, 0, len(opt.ReadDSN))
	for _, rd := range opt.ReadDSN {
		d, err := connect(opt, rd)
		if err != nil {
			return nil, err
		}
		addr = parseDSNAddr(rd)
		brk := brkGroup.Get(addr)
		r := &conn{DB: d, breaker: brk, opt: opt, addr: addr}
		rs = append(rs, r)
	}
	db.write = w
	db.read = rs
	db.master = &DB{write: db.write}
	return db, nil
}

func connect(opt *Options, dsn string) (*sql.DB, error) {
	d, err := sql.Open("mysql", dsn)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	d.SetMaxOpenConns(opt.Active)
	d.SetMaxIdleConns(opt.Idle)
	d.SetConnMaxLifetime(opt.IdleTimeout)
	return d, nil
}

// Begin starts a transaction. The isolation level is dependent on the driver.
func (db *DB) Begin(ctx context.Context) (tx *Tx, err error) {
	return db.write.begin(ctx)
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (db *DB) Exec(ctx context.Context, query string, args ...interface{}) (res sql.Result, err error) {
	return db.write.exec(ctx, query, args...)
}

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the returned
// statement. The caller must call the statement's Close method when the
// statement is no longer needed.
func (db *DB) Prepare(query string) (*Stmt, error) {
	return db.write.prepare(query)
}

// Prepared creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the returned
// statement. The caller must call the statement's Close method when the
// statement is no longer needed.
func (db *DB) Prepared(query string) (stmt *Stmt) {
	return db.write.prepared(query)
}

// Query executes a query that returns rows, typically a SELECT. The args are
// for any placeholder parameters in the query.
func (db *DB) Query(ctx context.Context, query string, args ...interface{}) (rows *Rows, err error) {
	idx := db.readIndex()
	for i := range db.read {
		rows, err = db.read[(idx+i)%len(db.read)].query(ctx, query, args...)
		if errs.EqualError(errs.ErrServiceUnavailable, err) {
			return
		}
	}
	return db.write.query(ctx, query, args...)
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until Row's
// Scan method is called.
func (db *DB) QueryRow(ctx context.Context, query string, args ...interface{}) *Row {
	idx := db.readIndex()
	for i := range db.read {
		row := db.read[(idx+i)%len(db.read)].queryRow(ctx, query, args...)
		if errs.EqualError(errs.ErrServiceUnavailable, row.err) {
			return row
		}
	}
	return db.write.queryRow(ctx, query, args...)
}

func (db *DB) readIndex() int {
	if len(db.read) == 0 {
		return 0
	}
	v := atomic.AddInt64(&db.idx, 1)
	return int(v) % len(db.read)
}

// Close closes the write and read database, releasing any open resources.
func (db *DB) Close() (err error) {
	if e := db.write.Close(); e != nil {
		err = errors.WithStack(e)
	}
	for _, rd := range db.read {
		if e := rd.Close(); e != nil {
			err = errors.WithStack(e)
		}
	}
	return
}

// Ping verifies a connection to the database is still alive, establishing a
// connection if necessary.
func (db *DB) Ping(ctx context.Context) (err error) {
	if err = db.write.ping(ctx); err != nil {
		return
	}
	for _, rd := range db.read {
		if err = rd.ping(ctx); err != nil {
			return
		}
	}
	return
}

// Master return *DB instance direct use master conn
// use this *DB instance only when you have some reason need to get result without any delay.
func (db *DB) Master() *DB {
	if db.master == nil {
		panic(ErrNoMaster)
	}
	return db.master
}

func (db *conn) onBreaker(err *error) {
	if err != nil && *err != nil && *err != sql.ErrNoRows && *err != sql.ErrTxDone {
		db.breaker.MarkFailed()
	} else {
		db.breaker.MarkSuccess()
	}
}

func (db *conn) begin(ctx context.Context) (tx *Tx, err error) {
	now := time.Now()
	defer slowLog("Begin", now)
	span, ctx := opentracing.StartSpanFromContext(ctx, "sql begin")
	ext.DBInstance.Set(span, db.addr)
	ext.DBType.Set(span, "sql")
	ext.SpanKind.Set(span, "client")
	defer func() {
		if err != nil {
			span.Finish()
		}
	}()
	if err = db.breaker.Allow(); err != nil {
		_metricReqErr.Inc(db.addr, db.addr, "begin", "breaker")
		return
	}
	ctx, cancel := shrink(ctx, db.opt.TranTimeout)
	rtx, err := db.BeginTx(ctx, nil)
	_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), db.addr, db.addr, "begin")
	if err != nil {
		err = errors.WithStack(err)
		cancel()
		return
	}
	tx = &Tx{tx: rtx, span: span, db: db, ctx: ctx, cancel: cancel}
	return
}

func (db *conn) exec(ctx context.Context, query string, args ...interface{}) (res sql.Result, err error) {
	now := time.Now()
	defer slowLog(fmt.Sprintf("Exec query(%s) args(%+v)", query, args), now)
	span, ctx := opentracing.StartSpanFromContext(ctx, "sql exec")
	ext.DBInstance.Set(span, db.addr)
	ext.DBStatement.Set(span, query)
	ext.DBType.Set(span, "sql")
	ext.SpanKind.Set(span, "client")
	defer span.Finish()
	if err = db.breaker.Allow(); err != nil {
		_metricReqErr.Inc(db.addr, db.addr, "exec", "breaker")
		return
	}
	ctx, cancel := shrink(ctx, db.opt.ExecTimeout)
	res, err = db.ExecContext(ctx, query, args...)
	cancel()
	db.onBreaker(&err)
	_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), db.addr, db.addr, "exec")
	if err != nil {
		err = errors.Wrapf(err, "exec:%s, args:%+v", query, args)
	}
	return
}

func (db *conn) ping(ctx context.Context) (err error) {
	now := time.Now()
	defer slowLog("Ping", now)
	span, ctx := opentracing.StartSpanFromContext(ctx, "sql ping")
	ext.DBInstance.Set(span, db.addr)
	ext.DBType.Set(span, "sql")
	ext.SpanKind.Set(span, "client")
	defer span.Finish()
	if err = db.breaker.Allow(); err != nil {
		_metricReqErr.Inc(db.addr, db.addr, "ping", "breaker")
		return
	}
	ctx, cancel := shrink(ctx, db.opt.ExecTimeout)
	err = db.PingContext(ctx)
	cancel()
	db.onBreaker(&err)
	_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), db.addr, db.addr, "ping")
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}

func (db *conn) prepare(query string) (*Stmt, error) {
	defer slowLog(fmt.Sprintf("Prepare query(%s)", query), time.Now())
	stmt, err := db.Prepare(query)
	if err != nil {
		err = errors.Wrapf(err, "prepare %s", query)
		return nil, err
	}
	st := &Stmt{query: query, db: db}
	st.stmt.Store(stmt)
	return st, nil
}

func (db *conn) prepared(query string) (stmt *Stmt) {
	defer slowLog(fmt.Sprintf("Prepared query(%s)", query), time.Now())
	stmt = &Stmt{query: query, db: db}
	s, err := db.Prepare(query)
	if err == nil {
		stmt.stmt.Store(s)
		return
	}
	go func() {
		for {
			s, err := db.Prepare(query)
			if err != nil {
				time.Sleep(time.Second)
				continue
			}
			stmt.stmt.Store(s)
			return
		}
	}()
	return
}

func (db *conn) query(ctx context.Context, query string, args ...interface{}) (rows *Rows, err error) {
	now := time.Now()
	defer slowLog(fmt.Sprintf("Query query(%s) args(%+v)", query, args), now)
	span, ctx := opentracing.StartSpanFromContext(ctx, "sql query")
	ext.DBInstance.Set(span, db.addr)
	ext.DBStatement.Set(span, query)
	ext.DBType.Set(span, "sql")
	ext.SpanKind.Set(span, "client")
	defer span.Finish()
	if err = db.breaker.Allow(); err != nil {
		_metricReqErr.Inc(db.addr, db.addr, "query", "breaker")
		return
	}
	ctx, cancel := shrink(ctx, db.opt.QueryTimeout)
	rs, err := db.DB.QueryContext(ctx, query, args...)
	db.onBreaker(&err)
	_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), db.addr, db.addr, "query")
	if err != nil {
		err = errors.Wrapf(err, "query:%s, args:%+v", query, args)
		cancel()
		return
	}
	rows = &Rows{Rows: rs, cancel: cancel}
	return
}

func (db *conn) queryRow(ctx context.Context, query string, args ...interface{}) *Row {
	now := time.Now()
	defer slowLog(fmt.Sprintf("QueryRow query(%s) args(%+v)", query, args), now)
	span, ctx := opentracing.StartSpanFromContext(ctx, "sql queryRow")
	ext.DBInstance.Set(span, db.addr)
	ext.DBStatement.Set(span, query)
	ext.DBType.Set(span, "sql")
	ext.SpanKind.Set(span, "client")
	if err := db.breaker.Allow(); err != nil {
		_metricReqErr.Inc(db.addr, db.addr, "queryRow", "breaker")
		return &Row{db: db, span: span, err: err}
	}
	ctx, cancel := shrink(ctx, db.opt.QueryTimeout)
	r := db.DB.QueryRowContext(ctx, query, args...)
	_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), db.addr, db.addr, "queryrow")
	return &Row{db: db, Row: r, query: query, args: args, span: span, cancel: cancel}
}

// Close closes the statement.
func (s *Stmt) Close() (err error) {
	if s == nil {
		err = ErrStmtNil
		return
	}
	stmt, ok := s.stmt.Load().(*sql.Stmt)
	if ok {
		err = errors.WithStack(stmt.Close())
	}
	return
}

// Exec executes a prepared statement with the given arguments and returns a
// Result summarizing the effect of the statement.
func (s *Stmt) Exec(ctx context.Context, args ...interface{}) (res sql.Result, err error) {
	if s == nil {
		err = ErrStmtNil
		return
	}
	now := time.Now()
	defer slowLog(fmt.Sprintf("Exec query(%s) args(%+v)", s.query, args), now)
	if s.tx {
		if s.span != nil {
			ext.DBStatement.Set(s.span, s.query)
		}
	} else {
		var span opentracing.Span
		span, ctx = opentracing.StartSpanFromContext(ctx, "sql exec")
		ext.DBInstance.Set(span, s.db.addr)
		ext.DBStatement.Set(span, s.query)
		ext.DBType.Set(span, "sql")
		ext.SpanKind.Set(span, "client")
		defer span.Finish()
	}
	if err = s.db.breaker.Allow(); err != nil {
		_metricReqErr.Inc(s.db.addr, s.db.addr, "stmt:exec", "breaker")
		return
	}
	stmt, ok := s.stmt.Load().(*sql.Stmt)
	if !ok {
		err = ErrStmtNil
		return
	}
	ctx, cancel := shrink(ctx, s.db.opt.ExecTimeout)
	res, err = stmt.ExecContext(ctx, args...)
	cancel()
	s.db.onBreaker(&err)
	_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), s.db.addr, s.db.addr, "stmt:exec")
	if err != nil {
		err = errors.Wrapf(err, "exec:%s, args:%+v", s.query, args)
	}
	return
}

// Query executes a prepared query statement with the given arguments and
// returns the query results as a *Rows.
func (s *Stmt) Query(ctx context.Context, args ...interface{}) (rows *Rows, err error) {
	if s == nil {
		err = ErrStmtNil
		return
	}
	now := time.Now()
	defer slowLog(fmt.Sprintf("Query query(%s) args(%+v)", s.query, args), now)
	if s.tx {
		if s.span != nil {
			ext.DBStatement.Set(s.span, s.query)
		}
	} else {
		var span opentracing.Span
		span, ctx = opentracing.StartSpanFromContext(ctx, "sql query")
		ext.DBInstance.Set(span, s.db.addr)
		ext.DBStatement.Set(span, s.query)
		ext.DBType.Set(span, "sql")
		ext.SpanKind.Set(span, "client")
		defer span.Finish()
	}
	if err = s.db.breaker.Allow(); err != nil {
		_metricReqErr.Inc(s.db.addr, s.db.addr, "stmt:query", "breaker")
		return
	}
	stmt, ok := s.stmt.Load().(*sql.Stmt)
	if !ok {
		err = ErrStmtNil
		return
	}
	ctx, cancel := shrink(ctx, s.db.opt.QueryTimeout)
	rs, err := stmt.QueryContext(ctx, args...)
	s.db.onBreaker(&err)
	_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), s.db.addr, s.db.addr, "stmt:query")
	if err != nil {
		err = errors.Wrapf(err, "query:%s, args:%+v", s.query, args)
		cancel()
		return
	}
	rows = &Rows{Rows: rs, cancel: cancel}
	return
}

// QueryRow executes a prepared query statement with the given arguments.
// If an error occurs during the execution of the statement, that error will
// be returned by a call to Scan on the returned *Row, which is always non-nil.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards the rest.
func (s *Stmt) QueryRow(ctx context.Context, args ...interface{}) (row *Row) {
	now := time.Now()
	defer slowLog(fmt.Sprintf("QueryRow query(%s) args(%+v)", s.query, args), now)
	row = &Row{db: s.db, query: s.query, args: args}
	if s == nil {
		row.err = ErrStmtNil
		return
	}
	if s.tx {
		if s.span != nil {
			ext.DBStatement.Set(s.span, s.query)
		}
	} else {
		var span opentracing.Span
		span, ctx = opentracing.StartSpanFromContext(ctx, "sql queryRow")
		ext.DBInstance.Set(span, s.db.addr)
		ext.DBStatement.Set(span, s.query)
		ext.DBType.Set(span, "sql")
		ext.SpanKind.Set(span, "client")
		defer span.Finish()
	}
	if row.err = s.db.breaker.Allow(); row.err != nil {
		_metricReqErr.Inc(s.db.addr, s.db.addr, "stmt:queryrow", "breaker")
		return
	}
	stmt, ok := s.stmt.Load().(*sql.Stmt)
	if !ok {
		return
	}
	ctx, cancel := shrink(ctx, s.db.opt.QueryTimeout)
	row.Row = stmt.QueryRowContext(ctx, args...)
	row.cancel = cancel
	_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), s.db.addr, s.db.addr, "stmt:queryrow")
	return
}

// Commit commits the transaction.
func (tx *Tx) Commit() (err error) {
	err = tx.tx.Commit()
	tx.cancel()
	tx.db.onBreaker(&err)
	if tx.span != nil {
		tx.span.Finish()
	}
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}

// Rollback aborts the transaction.
func (tx *Tx) Rollback() (err error) {
	err = tx.tx.Rollback()
	tx.cancel()
	tx.db.onBreaker(&err)
	if tx.span != nil {
		tx.span.Finish()
	}
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}

// Exec executes a query that doesn't return rows. For example: an INSERT and
// UPDATE.
func (tx *Tx) Exec(query string, args ...interface{}) (res sql.Result, err error) {
	now := time.Now()
	defer slowLog(fmt.Sprintf("Exec query(%s) args(%+v)", query, args), now)
	if tx.span != nil {
		// TODO args 参数有可能比较大，先不放入 trace，后面有需要的话可以考虑加入，其他trace地方相同处理
		//tx.span.SetTag("db.params", fmt.Sprintf("%v", args))
		ext.DBStatement.Set(tx.span, fmt.Sprintf("exec %s", query))
	}
	res, err = tx.tx.ExecContext(tx.ctx, query, args...)
	_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), tx.db.addr, tx.db.addr, "tx:exec")
	if err != nil {
		err = errors.Wrapf(err, "exec:%s, args:%+v", query, args)
	}
	return
}

// Query executes a query that returns rows, typically a SELECT.
func (tx *Tx) Query(query string, args ...interface{}) (rows *Rows, err error) {
	if tx.span != nil {
		ext.DBStatement.Set(tx.span, fmt.Sprintf("query %s", query))
	}
	now := time.Now()
	defer slowLog(fmt.Sprintf("Query query(%s) args(%+v)", query, args), now)
	defer func() {
		_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), tx.db.addr, tx.db.addr, "tx:query")
	}()
	rs, err := tx.tx.QueryContext(tx.ctx, query, args...)
	if err == nil {
		rows = &Rows{Rows: rs}
	} else {
		err = errors.Wrapf(err, "query:%s, args:%+v", query, args)
	}
	return
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until Row's
// Scan method is called.
func (tx *Tx) QueryRow(query string, args ...interface{}) *Row {
	if tx.span != nil {
		ext.DBStatement.Set(tx.span, fmt.Sprintf("queryRow %s", query))
	}
	now := time.Now()
	defer slowLog(fmt.Sprintf("QueryRow query(%s) args(%+v)", query, args), now)
	defer func() {
		_metricReqDur.Observe(int64(time.Since(now)/time.Millisecond), tx.db.addr, tx.db.addr, "tx:queryrow")
	}()
	r := tx.tx.QueryRowContext(tx.ctx, query, args...)
	return &Row{Row: r, db: tx.db, query: query, args: args}
}

// Stmt returns a transaction-specific prepared statement from an existing statement.
func (tx *Tx) Stmt(stmt *Stmt) *Stmt {
	as, ok := stmt.stmt.Load().(*sql.Stmt)
	if !ok {
		return nil
	}
	ts := tx.tx.StmtContext(tx.ctx, as)
	st := &Stmt{query: stmt.query, tx: true, span: tx.span, db: tx.db}
	st.stmt.Store(ts)
	return st
}

// Prepare creates a prepared statement for use within a transaction.
// The returned statement operates within the transaction and can no longer be
// used once the transaction has been committed or rolled back.
// To use an existing prepared statement on this transaction, see Tx.Stmt.
func (tx *Tx) Prepare(query string) (*Stmt, error) {
	if tx.span != nil {
		ext.DBStatement.Set(tx.span, fmt.Sprintf("prepare %s", query))
	}
	defer slowLog(fmt.Sprintf("Prepare query(%s)", query), time.Now())
	stmt, err := tx.tx.Prepare(query)
	if err != nil {
		err = errors.Wrapf(err, "prepare %s", query)
		return nil, err
	}
	st := &Stmt{query: query, tx: true, span: tx.span, db: tx.db}
	st.stmt.Store(stmt)
	return st, nil
}

// parseDSNAddr parse dsn name and return addr.
func parseDSNAddr(dsn string) (addr string) {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		// just ignore parseDSN error, mysql client will return error for us when connect.
		return ""
	}
	return cfg.Addr
}

func slowLog(statement string, now time.Time) {
	du := time.Since(now)
	if du > _slowLogDuration {
		log.Warnf("%s slow log statement: %s time: %v", _family, statement, du)
	}
}

// shrink will decrease the duration by comparing with context's timeout duration
// and return new timeout\context\CancelFunc.
func shrink(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	if deadline, ok := ctx.Deadline(); ok {
		if tm := time.Until(deadline); tm < d {
			// deliver small timeout
			return ctx, func() {}
		}
	}
	ctx, cancel := context.WithTimeout(ctx, d)
	return ctx, cancel
}
