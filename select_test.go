package mapper

import "testing"

type selectUser struct {
	_        interface{} `table:"user"`
	Id       int64       `column:"id"`
	Username string      `column:"username"`
}

func TestDbSelect_toSql(t *testing.T) {
	AddStruct(selectUser{})

	t.Run("test pointer", func(t *testing.T) {
		DbSelectToSql(t, &selectUser{
			Id:       1,
			Username: "John",
		})
	})

	t.Run("test empty", func(t *testing.T) {
		DbSelectToSql(t, selectUser{})
	})

	t.Run("test value", func(t *testing.T) {
		DbSelectToSql(t, selectUser{
			Id:       1,
			Username: "John",
		})
	})
}

func DbSelectToSql(t *testing.T, c interface{}) {
	s, err := toSelect(c)
	if err != nil {
		t.Skip(err.Error())
		return
	}

	info := s.toSql()
	t.Log(info)
}
