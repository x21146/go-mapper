package mapper

import (
	"fmt"
	"strings"
)

type dbInsertNotExists struct {
	*dbInsert
	where map[string]interface{}
}

func toInsertNotExists(d, c interface{}) (*dbInsertNotExists, error) {
	i, err := toInsert(d)
	if err != nil {
		return nil, err
	}

	return &dbInsertNotExists{
		dbInsert: i,
		where:    where(c, i.info),
	}, nil
}

func (i *dbInsertNotExists) toSql() *sqlInfo {
	info := &sqlInfo{}

	var cs, vs []string
	for column, value := range i.values {
		cs = append(cs, fmt.Sprintf("`%s`", column))
		vs = append(vs, "?")
		info.args = append(info.args, value)
	}

	var ws []string
	for column, value := range i.where {
		ws = append(ws, fmt.Sprintf("`%s` = ?", column))
		info.args = append(info.args, value)
	}

	notExists := ""
	if len(ws) > 0 {
		condition := strings.Join(ws, " AND ")
		notExists = fmt.Sprintf(" WHERE NOT EXISTS (SELECT 1 FROM `%s` WHERE %s)", i.tableName, condition)
	}

	columns := strings.Join(cs, ", ")
	values := strings.Join(vs, ", ")
	info.sql = fmt.Sprintf("INSERT INTO `%s` (%s) SELECT %s%s", i.tableName, columns, values, notExists)

	return info
}
