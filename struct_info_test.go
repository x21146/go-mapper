package mapper

import (
	"testing"
)

type AddStructTest struct {
	_  interface{} `table:"user"`
	Id int64       `column:"id"`
}

func TestToStructInfo(t *testing.T) {
	t.Run("parse struct", func(t *testing.T) {
		if i, err := toStructInfo(AddStructTest{}); err != nil || i == nil {
			t.FailNow()
		}
	})

	t.Run("parse pointer", func(t *testing.T) {
		if i, err := toStructInfo(&AddStructTest{}); err != nil || i == nil {
			t.FailNow()
		}
	})

	t.Run("parse nil", func(t *testing.T) {
		var n *AddStructTest = nil
		if i, err := toStructInfo(n); err != nil || i == nil {
			t.FailNow()
		}
	})
}
