package go_mapper

import "reflect"

type structInfo struct {
	name   string
	table  string
	fields map[string]reflect.StructField
}

var emptyValue = reflect.Value{}

var structs = make(map[string]*structInfo)

func getInfo(name string) *structInfo {
	if info, ok := structs[name]; ok {
		return info
	} else {
		return nil
	}
}

func checkAndGetInfo(d interface{}) (reflect.Value, *structInfo, error) {
	if d == nil {
		return emptyValue, nil, ErrNilData
	}

	v := reflect.ValueOf(d)
	// get pointer value
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	info := getInfo(v.Type().Name())
	if info == nil {
		return emptyValue, nil, ErrStructInfoNotExists
	}

	if v.Interface() == reflect.Zero(v.Type()).Interface() {
		return v, info, ErrEmptyData
	}

	return v, info, nil
}

func AddStruct(s interface{}) {
	info, err := toStructInfo(s)
	if err != nil {
		return
	}

	structs[info.name] = info
}

func toStructInfo(s interface{}) (info *structInfo, err error) {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}

	info = &structInfo{
		name:   t.Name(),
		fields: make(map[string]reflect.StructField),
	}

	n := t.NumField()
	for i := 0; i < n; i++ {
		f := t.Field(i)

		if table, ok := f.Tag.Lookup("table"); ok {
			info.table = table
			continue
		}

		if column, ok := f.Tag.Lookup("column"); ok {
			info.fields[column] = f
			continue
		}
	}

	if info.name == "" {
		return nil, ErrAnonymousStruct
	}

	if len(info.fields) == 0 {
		return nil, ErrNoColumnInStruct
	}

	if info.table == "" {
		return nil, ErrNoTableInStruct
	}

	return info, nil
}
