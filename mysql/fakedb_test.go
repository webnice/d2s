package mysql

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"sync"
	"testing"
)

// Скриптованный драйвер database/sql для тестов без внешних зависимостей.
type fakeDriver struct{}
type fakeConn struct{}
type fakeStmt struct{ query string }
type fakeTx struct{}
type fakeRows struct {
	columns  []string
	values   [][]driver.Value
	pos      int
	closeErr error
}

var (
	fakeMu       sync.Mutex
	fakeQuery    func(query string) (columns []string, values [][]driver.Value, err error)
	fakeCloseErr error
)

func (d *fakeDriver) Open(name string) (driver.Conn, error) { return &fakeConn{}, nil }
func (c *fakeConn) Prepare(query string) (driver.Stmt, error) {
	return &fakeStmt{query: query}, nil
}
func (c *fakeConn) Close() error              { return nil }
func (c *fakeConn) Begin() (driver.Tx, error) { return &fakeTx{}, nil }
func (t *fakeTx) Commit() error               { return nil }
func (t *fakeTx) Rollback() error             { return nil }
func (s *fakeStmt) Close() error              { return nil }
func (s *fakeStmt) NumInput() int             { return -1 }
func (s *fakeStmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, errors.New("не реализовано")
}
func (s *fakeStmt) Query(args []driver.Value) (driver.Rows, error) {
	fakeMu.Lock()
	defer fakeMu.Unlock()
	if fakeQuery == nil {
		return nil, errors.New("скрипт запроса не установлен")
	}
	cols, vals, err := fakeQuery(s.query)
	if err != nil {
		return nil, err
	}
	return &fakeRows{columns: cols, values: vals, closeErr: fakeCloseErr}, nil
}
func (r *fakeRows) Columns() []string { return r.columns }
func (r *fakeRows) Close() error      { return r.closeErr }
func (r *fakeRows) Next(dest []driver.Value) error {
	if r.pos >= len(r.values) {
		return io.EOF
	}
	copy(dest, r.values[r.pos])
	r.pos++

	return nil
}

func init() { sql.Register("d2s-fake-mysql", &fakeDriver{}) }

// fakeDB Возвращает скриптованное соединение с базой данных.
func fakeDB(t *testing.T, query func(string) ([]string, [][]driver.Value, error)) *sql.DB {
	t.Helper()
	fakeMu.Lock()
	fakeQuery = query
	fakeCloseErr = nil
	fakeMu.Unlock()
	db, err := sql.Open("d2s-fake-mysql", "")
	if err != nil {
		t.Fatal(err)
	}

	return db
}

// fakeDBCloseErr Возвращает соединение, у которого rows.Close() завершается ошибкой.
func fakeDBCloseErr(t *testing.T, query func(string) ([]string, [][]driver.Value, error)) *sql.DB {
	t.Helper()
	db := fakeDB(t, query)
	fakeMu.Lock()
	fakeCloseErr = errors.New("ошибка закрытия строк")
	fakeMu.Unlock()

	return db
}
