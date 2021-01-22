package mapper

import (
	"fmt"
	"strings"
)

type dbInsertUpdate struct {
	*dbInsert
}

func toInsertUpdate(d interface{}) (*dbInsertUpdate, error) {
	i, err := toInsert(d)
	if err != nil {
		return nil, err
	}

	return &dbInsertUpdate{dbInsert: i}, nil
}

func (i *dbInsertUpdate) toSql() *sqlInfo {
	info := i.dbInsert.toSql()

	var cs []string
	for col, _ := range i.values {
		cs = append(cs, fmt.Sprintf("`%s` = VALUES(`%s`)", col, col))
	}

	cols := strings.Join(cs, ", ")
	info.sql += fmt.Sprintf(" ON DUPLICATE KEY UPDATE %s", cols)

	return info
}
