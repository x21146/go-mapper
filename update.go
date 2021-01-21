package go_mapper

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type dbUpdate struct {
	tableName string
	values    map[string]interface{}
	args      map[string]interface{}
}

func where(c interface{}, info *structInfo) (args map[string]interface{}) {
	args = make(map[string]interface{})

	v := reflect.ValueOf(c)
	isNil := v.Interface() == reflect.Zero(v.Type())

	// get pointer value
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if !isNil {
		for col, f := range info.fields {
			fv := v.FieldByName(f.Name)
			vi := fv.Interface()
			if vi == reflect.Zero(f.Type).Interface() {
				// skip zero value in dbSelect condition
				continue
			}

			args[col] = vi
		}
	}

	return
}

func toUpdate(d, c interface{}) (*dbUpdate, error) {
	v, info, err := checkAndGetInfo(d)
	if err != nil {
		return nil, err
	}

	u := &dbUpdate{
		tableName: info.table,
		values:    make(map[string]interface{}),
		args:      where(c, info),
	}

	for col, f := range info.fields {
		fv := v.FieldByName(f.Name)
		vi := fv.Interface()
		if vi == reflect.Zero(f.Type).Interface() {
			// skip zero value in dbSelect condition
			continue
		}

		u.values[col] = fv.Interface()
	}

	return u, nil
}

func (u *dbUpdate) toSql() *sqlInfo {
	info := &sqlInfo{}

	var values []string
	for column, value := range u.values {
		values = append(values, fmt.Sprintf("`%s` = ?", column))
		info.args = append(info.args, value)
	}

	set := ""
	if len(values) > 0 {
		set = strings.Join(values, ", ")
	}

	var condition []string
	for column, value := range u.args {
		condition = append(condition, fmt.Sprintf("`%s` = ?", column))
		info.args = append(info.args, value)
	}

	where := ""
	if len(condition) > 0 {
		where = " WHERE " + strings.Join(condition, " AND ")
	}

	info.sql = fmt.Sprintf("UPDATE %s SET %s%s", u.tableName, set, where)
	return info
}

func (u *dbUpdate) Do(db *sql.DB) (int64, error) {
	info := u.toSql()
	res, err := info.exec(db)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
