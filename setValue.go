package libconfig

import (
	"reflect"
	"strconv"
)

// setValue parses the bytes into a reflect.Value
func setValue(v reflect.Value, key string, value []byte) error {
	var f func(reflect.Value, reflect.Kind, string, string) error
	k := v.Kind()

	switch k {

	// []byte
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			v.SetBytes(value)
			return nil
		}

	// string
	case reflect.String:
		v.SetString(string(value))
		return nil

	// int
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		f = setValueToInt

	// uint
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		f = setValueToUint

	// float
	case reflect.Float32, reflect.Float64:
		f = setValueToFloat

	// bool
	case reflect.Bool:
		f = setValueToBool
	}

	if f == nil {
		return NewErrCannotSetKind(k)
	}

	return f(v, k, key, string(value))
}

func setValueToInt(v reflect.Value, k reflect.Kind, key, value string) error {
	intVal, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return NewErrCannotParseEnv(err, k, key, value)
	}

	if v.OverflowInt(intVal) {
		return NewErrOverflow(k, key, value)
	}

	v.SetInt(intVal)
	return nil
}

func setValueToUint(v reflect.Value, k reflect.Kind, key, value string) error {
	uintVal, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return NewErrCannotParseEnv(err, k, key, value)
	}

	if v.OverflowUint(uintVal) {
		return NewErrOverflow(k, key, value)
	}

	v.SetUint(uintVal)
	return nil
}

func setValueToFloat(v reflect.Value, k reflect.Kind, key, value string) error {
	floatVal, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return NewErrCannotParseEnv(err, k, key, value)
	}

	if v.OverflowFloat(floatVal) {
		return NewErrOverflow(k, key, value)
	}

	v.SetFloat(floatVal)
	return nil
}

func setValueToBool(v reflect.Value, k reflect.Kind, key, value string) error {
	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		return NewErrCannotParseEnv(err, k, key, value)
	}

	v.SetBool(boolVal)
	return nil
}
