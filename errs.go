package go_mapper

import (
	"errors"
)

var (
	ErrNilData             = errors.New("data is nil for sql")
	ErrEmptyData           = errors.New("data is empty for sql")
	ErrNotStruct           = errors.New("data is not struct type")
	ErrStructInfoNotExists = errors.New("struct info not exists")
	ErrAnonymousStruct     = errors.New("not support anonymous struct")
	ErrNoColumnInStruct    = errors.New("not found fields in struct using 'column' tag")
	ErrNoTableInStruct     = errors.New("not found table in struct using 'table' tag")
	ErrScanNotPointer      = errors.New("scan data need pointer output")
)
