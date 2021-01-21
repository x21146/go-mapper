package go_mapper

import (
	"database/sql"
	"reflect"
)

type scanner func(*sql.Rows, reflect.Value) error

func scan(row *sql.Rows, out interface{}) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr {
		return ErrScanNotPointer
	}

	v = v.Elem()

	var s scanner
	switch v.Type().Kind() {
	case reflect.Array, reflect.Slice:
		s = scanSlice
	default:
		s = scanStruct
	}

	return s(row, v)
}

func scanStruct(row *sql.Rows, v reflect.Value) error {
	defer row.Close()

	info := getInfo(v.Type().Name())

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

func scanSlice(row *sql.Rows, v reflect.Value) error {
	v.Set(reflect.MakeSlice(v.Type(), 0, 0))
	et := v.Type().Elem()

	cols, err := row.Columns()
	if err != nil {
		return err
	}

	info := getInfo(et.Name())
	if info == nil {
		return ErrStructInfoNotExists
	}

	defer row.Close()
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
