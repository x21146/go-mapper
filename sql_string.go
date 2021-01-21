package go_mapper

import "database/sql"

type sqlInfo struct {
	sql  string
	args []interface{}
}

func (i *sqlInfo) query(db *sql.DB) (*sql.Rows, error) {
	return db.Query(i.sql, i.args...)
}

func (i *sqlInfo) queryRow(db *sql.DB) *sql.Row {
	return db.QueryRow(i.sql, i.args...)
}

func (i *sqlInfo) exec(db *sql.DB) (sql.Result, error) {
	return db.Exec(i.sql, i.args...)
}
