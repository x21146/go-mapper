package go_mapper

import "testing"

type insertUser struct {
	_        interface{} `table:"user"`
	Id       int64       `column:"id"`
	Username string      `column:"username"`
}

func TestDbInsert_toSql(t *testing.T) {
	var s *insertUser = nil
	AddStruct(s)

	t.Run("test pointer", func(t *testing.T) {
		dbInsertToSql(t, &insertUser{Username: "John"})
	})

	t.Run("test value", func(t *testing.T) {
		dbInsertToSql(t, insertUser{Username: "John"})
	})
}

func dbInsertToSql(t *testing.T, d interface{}) {
	i, err := toInsert(d)
	if err != nil {
		t.Skip(err.Error())
	}

	info := i.toSql()
	t.Log(info)
}
