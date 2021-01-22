package mapper

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

func (m *BaseMapper) do(d doResult, err error) (int64, error) {
	if err != nil {
		return 0, err
	}

	return d.Do(m.Db)
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
	return m.do(toInsert(data))
}

func (m *BaseMapper) InsertUpdate(data interface{}) (int64, error) {
	return m.do(toInsertUpdate(data))
}

func (m *BaseMapper) InsertNotExists(data, condition interface{}) (int64, error) {
	return m.do(toInsertNotExists(data, condition))
}

func (m *BaseMapper) Update(data, condition interface{}) (int64, error) {
	return m.do(toUpdate(data, condition))
}

func (m *BaseMapper) Delete(condition interface{}) (int64, error) {
	return m.do(toDelete(condition))
}
