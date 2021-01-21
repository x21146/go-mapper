package go_mapper

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"testing"
	"time"
)

type testStruct struct {
	_        interface{} `table:"test_table"`
	Id       int64       `column:"id"`
	Username string      `column:"username"`
	Password string      `column:"password"`
}

type testMapper struct {
	BaseMapper
}

const createTable = "CREATE TABLE IF NOT EXISTS `test_table` (\n`id` bigint NOT NULL AUTO_INCREMENT,\n`username` varchar(255) NOT NULL,\n`password` varchar(255) NOT NULL,\nPRIMARY KEY (`id`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"
const dropTable = "DROP TABLE IF EXISTS `test_table`;"

func (m *testMapper) selectById(id int64) (*testStruct, error) {
	data := &testStruct{}
	return data, m.Select(testStruct{Id: id}, data)
}

func mockDb(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/testing")
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
		return nil
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}

func TestBaseMapper(t *testing.T) {
	var s *testStruct = nil
	AddStruct(s)

	var id1 *int64
	var id2 *int64
	wg := &sync.WaitGroup{}

	mapper := &testMapper{BaseMapper{Db: mockDb(t)}}
	_, _ = mapper.Db.Exec(createTable)

	wg.Add(1)
	t.Run("insert 1", func(t *testing.T) {
		defer wg.Done()
		id, err := mapper.Insert(testStruct{Username: "John", Password: "password1"})
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}

		id1 = &id
		t.Log("insert first data, id:", id)
	})

	wg.Add(1)
	t.Run("insert 2", func(t *testing.T) {
		defer wg.Done()
		id, err := mapper.Insert(testStruct{Username: "Mike", Password: "password2"})
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}

		id2 = &id
		t.Log("insert second data, id:", id)
	})

	wg.Add(1)
	t.Run("update", func(t *testing.T) {
		defer wg.Done()
		r, err := mapper.Update(testStruct{Password: "newPassword"}, testStruct{Id: *id2})
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}

		if r == 1 {
			t.Log("update new password")
		}
	})

	wg.Add(1)
	t.Run("single select test", func(t *testing.T) {
		defer wg.Done()
		data := &testStruct{}
		err := mapper.Select(testStruct{Id: *id1}, data)
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}

		t.Log("select data at", *id1, ":", data)
	})

	wg.Add(1)
	t.Run("select slice test", func(t *testing.T) {
		defer wg.Done()
		var data []testStruct
		err := mapper.Select(testStruct{}, &data)
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}

		t.Log("select data slice:", data)
	})

	wg.Wait()
	_, _ = mapper.Db.Exec(dropTable)
}
