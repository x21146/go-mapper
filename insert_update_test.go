package mapper

import "testing"

func TestDbInsertUpdate_toSql(t *testing.T) {
	var s *insertUser = nil
	AddStruct(s)

	t.Run("test pointer", func(t *testing.T) {
		dbInsertUpdateToSql(t, &insertUser{Username: "John"})
	})

	t.Run("test value", func(t *testing.T) {
		dbInsertUpdateToSql(t, insertUser{Username: "John"})
	})
}

func dbInsertUpdateToSql(t *testing.T, d interface{}) {
	i, err := toInsertUpdate(d)
	if err != nil {
		t.Skip(err.Error())
	}

	info := i.toSql()
	t.Log(info)
}
