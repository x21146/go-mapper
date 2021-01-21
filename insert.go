package go_mapper

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type dbInsert struct {
	tableName string
	values    map[string]interface{}
}

func toInsert(d interface{}) (*dbInsert, error) {
	v, info, err := checkAndGetInfo(d)
	if err != nil {
		return nil, err
	}

	i := &dbInsert{
		tableName: info.table,
		values:    make(map[string]interface{}),
	}

	for col, f := range info.fields {
		fv := v.FieldByName(f.Name)
		vi := fv.Interface()
		if vi == reflect.Zero(f.Type).Interface() {
			// skip zero value in dbSelect condition
			continue
		}

		i.values[col] = fv.Interface()
	}

	return i, nil
}

func (i *dbInsert) toSql() *sqlInfo {
	info := &sqlInfo{}

	var cs, vs []string
	for column, value := range i.values {
		cs = append(cs, fmt.Sprintf("`%s`", column))
		vs = append(vs, "?")
		info.args = append(info.args, value)
	}

	columns := strings.Join(cs, ", ")
	values := strings.Join(vs, ", ")
	info.sql = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", i.tableName, columns, values)

	return info
}

func (i *dbInsert) Do(db *sql.DB) (int64, error) {
	info := i.toSql()
	res, err := info.exec(db)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}
