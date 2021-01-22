package mapper

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type dbDelete struct {
	tableName string
	args      map[string]interface{}
}

func toDelete(c interface{}) (*dbDelete, error) {
	v, info, err := checkAndGetInfo(c)
	if err != nil {
		return nil, err
	}
	if err != nil && err != ErrEmptyData {
		return nil, err
	}

	if info == nil {
		return nil, ErrStructInfoNotExists
	}

	isNil := err == ErrEmptyData

	d := &dbDelete{
		tableName: info.table,
		args:      make(map[string]interface{}),
	}

	if !isNil {
		for col, f := range info.fields {
			fv := v.FieldByName(f.Name)
			vi := fv.Interface()
			if vi == reflect.Zero(f.Type).Interface() {
				// skip zero value in dbSelect condition
				continue
			}

			d.args[col] = vi
		}
	}

	return d, nil
}

func (d *dbDelete) toSql() *sqlInfo {
	info := &sqlInfo{}

	var condition []string
	for column, value := range d.args {
		condition = append(condition, fmt.Sprintf("`%s` = ?", column))
		info.args = append(info.args, value)
	}

	where := ""
	if len(condition) > 0 {
		where = " WHERE " + strings.Join(condition, " AND ")
	}

	info.sql = fmt.Sprintf("DELETE * FROM `%s`%s", d.tableName, where)
	return info
}

func (d *dbDelete) Do(db *sql.DB) (int64, error) {
	info := d.toSql()
	res, err := info.exec(db)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
