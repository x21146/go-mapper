package go_mapper

import (
	"database/sql"
	"log"
)

type BaseMapper struct {
	Db     *sql.DB
	logger *log.Logger
}

func (m *BaseMapper) log(args ...interface{}) {
	if m.logger == nil {
		log.Println(args...)
	} else {
		m.logger.Println(args...)
	}
}

func (m *BaseMapper) Select(condition, out interface{}) error {
	s, err := toSelect(condition)
	if err != nil {
		return err
	}

	return s.Do(m.Db, out)
}

func (m *BaseMapper) SelectCount(condition interface{}) (int64, error) {
	s, err := toSelect(condition)
	if err != nil {
		return 0, err
	}

	return s.DoCount(m.Db)
}

func (m *BaseMapper) Insert(data interface{}) (int64, error) {
	i, err := toInsert(data)
	if err != nil {
		return 0, err
	}

	return i.Do(m.Db)
}

func (m *BaseMapper) Update(data, condition interface{}) (int64, error) {
	u, err := toUpdate(data, condition)
	if err != nil {
		return 0, err
	}

	return u.Do(m.Db)
}
