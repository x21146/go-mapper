package mapper

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type dbSelect struct {
	tableName string
	columns   []string
	args      map[string]interface{}
}

func toSelect(c interface{}) (*dbSelect, error) {
	v, info, err := checkAndGetInfo(c)
	if err != nil && err != ErrEmptyData {
		return nil, err
	}

	if info == nil {
		return nil, ErrStructInfoNotExists
	}

	isNil := err == ErrEmptyData

	s := &dbSelect{
		tableName: info.table,
		args:      make(map[string]interface{}),
	}

	for col, f := range info.fields {
		s.columns = append(s.columns, fmt.Sprintf("`%s`", col))

		if !isNil {
			fv := v.FieldByName(f.Name)
			vi := fv.Interface()
			if vi == reflect.Zero(f.Type).Interface() {
				// skip zero value in dbSelect condition
				continue
			}

			s.args[col] = vi
		}
	}

	return s, nil
}

func (s *dbSelect) toSql() *sqlInfo {
	info := &sqlInfo{}

	columns := strings.Join(s.columns, ", ")
	var condition []string
	for column, value := range s.args {
		condition = append(condition, fmt.Sprintf("`%s` = ?", column))
		info.args = append(info.args, value)
	}

	where := ""
	if len(condition) > 0 {
		where = " WHERE " + strings.Join(condition, " AND ")
	}

	info.sql = fmt.Sprintf("SELECT `%s` FROM %s%s", columns, s.tableName, where)
	return info
}

func (s *dbSelect) toCountSql() *sqlInfo {
	info := &sqlInfo{}

	var condition []string
	for column, value := range s.args {
		condition = append(condition, fmt.Sprintf("`%s` = ?", column))
		info.args = append(info.args, value)
	}

	where := ""
	if len(condition) > 0 {
		where = " WHERE " + strings.Join(condition, " AND ")
	}

	info.sql = fmt.Sprintf("SELECT COUNT(0) FROM %s%s", s.tableName, where)
	return info
}

func (s *dbSelect) Do(db *sql.DB, out interface{}) error {
	info := s.toSql()
	row, err := info.query(db)
	if err != nil {
		return err
	}

	return scan(row, out)
}

func (s *dbSelect) DoCount(db *sql.DB) (int64, error) {
	info := s.toCountSql()
	row := info.queryRow(db)

	var c int64 = 0
	return c, row.Scan(&c)
}
