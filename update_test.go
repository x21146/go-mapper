package mapper

import "testing"

type updateUser struct {
	_        interface{} `table:"user"`
	Id       int64       `column:"id"`
	Username string      `column:"username"`
}

func TestDbUpdate_toSql(t *testing.T) {
	AddStruct(updateUser{})

	t.Run("test pointer", func(t *testing.T) {
		dbUpdateToSql(t, &updateUser{
			Id:       1,
			Username: "John",
		}, &updateUser{
			Id: 1,
		})
	})

	t.Run("test empty", func(t *testing.T) {
		dbUpdateToSql(t, updateUser{
			Username: "John",
		}, updateUser{})
	})

	t.Run("test value", func(t *testing.T) {
		dbUpdateToSql(t, updateUser{
			Id:       1,
			Username: "John",
		}, updateUser{
			Id: 1,
		})
	})
}

func dbUpdateToSql(t *testing.T, d, c interface{}) {
	u, err := toUpdate(d, c)
	if err != nil {
		t.Skip(err.Error())
	}

	info := u.toSql()
	t.Log(info)
}
