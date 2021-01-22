package mapper

import "testing"

func TestDbInsertNotExists_ToSql(t *testing.T) {
	var s *insertUser = nil
	AddStruct(s)

	t.Run("test pointer", func(t *testing.T) {
		dbInsertNotExistsToSql(t, &insertUser{Username: "John"}, insertUser{Id: 1})
	})

	t.Run("test value", func(t *testing.T) {
		dbInsertNotExistsToSql(t, insertUser{Username: "John"}, insertUser{Id: 1})
	})
}

func dbInsertNotExistsToSql(t *testing.T, d, c interface{}) {
	i, err := toInsertNotExists(d, c)
	if err != nil {
		t.Skip(err.Error())
	}

	info := i.toSql()
	t.Log(info)
}
