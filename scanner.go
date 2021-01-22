package mapper

import (
	"database/sql"
	"reflect"
)

type scanner func(*sql.Rows, reflect.Value, reflect.Type) error

func scan(row *sql.Rows, out interface{}) error {
	//v := reflect.ValueOf(out)
	//if v.Kind() != reflect.Ptr {
	//	return ErrScanNotPointer
	//}
	//
	//v = v.Elem()
	//
	//var s scanner
	//switch v.Type().Kind() {
	//case reflect.Array, reflect.Slice:
	//	s = scanSlice
	//default:
	//	s = scanStruct
	//}
	//
	//return s(row, v)

	t := reflect.TypeOf(out)
	if t.Kind() != reflect.Ptr {
		return ErrScanNotPointer
	}

	t = t.Elem()

	var s scanner
	switch t.Kind() {
	case reflect.Array, reflect.Slice:
		s = scanSlice
	default:
		s = scanStruct
	}

	return s(row, reflect.ValueOf(out), t)
}

func scanStruct(row *sql.Rows, v reflect.Value, t reflect.Type) error {
	defer row.Close()

	v = v.Elem()
	info := getInfo(t.Name())

	if !row.Next() {
		return nil
	}

	cols, err := row.Columns()
	if err != nil {
		return err
	}

	var dest []interface{}
	for _, col := range cols {
		f, ok := info.fields[col]
		if !ok {
			continue
		}

		vf := v.FieldByName(f.Name)

		var i interface{}
		if vf.CanAddr() {
			i = vf.Addr().Interface()
		} else {
			i = vf.Interface()
		}

		dest = append(dest, i)
	}

	return row.Scan(dest...)
}

func scanSlice(row *sql.Rows, v reflect.Value, t reflect.Type) error {
	defer row.Close()

	v = v.Elem()
	v.Set(reflect.MakeSlice(t, 0, 0))
	et := v.Type().Elem()

	cols, err := row.Columns()
	if err != nil {
		return err
	}

	info := getInfo(et.Name())
	if info == nil {
		return ErrStructInfoNotExists
	}

	for row.Next() {
		var dest []interface{}
		evp := reflect.New(et)
		ev := evp.Elem()
		for _, col := range cols {
			f, ok := info.fields[col]
			if !ok {
				continue
			}

			evf := ev.FieldByName(f.Name)

			var i interface{}
			if evf.CanAddr() {
				i = evf.Addr().Interface()
			} else {
				i = evf.Interface()
			}

			dest = append(dest, i)
		}

		err := row.Scan(dest...)
		if err != nil {
			_ = row.Close()
			return nil
		}

		v.Set(reflect.Append(v, ev))
	}

	return nil
}
