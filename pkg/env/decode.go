package env

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

// Decode environment variable to given v.
func Decode[T any](v *T) {
	if v == nil {
		panic(errors.New("arg of Decode must not be nil pointer"))
	}
	rv := reflect.ValueOf(v).Elem()
	for i := 0; i < rv.NumField(); i++ {
		rf := rv.Type().Field(i)
		tag := rf.Tag.Get("env")
		if tag == "" {
			continue
		}
		if v, ok := os.LookupEnv(tag); ok {
			rv.Field(i).Set(assign(v, rf.Type))
		}
	}
}

func assign(v string, typ reflect.Type) reflect.Value {
	switch typ.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(v)
		if err != nil {
			panic(err)
		}
		return reflect.ValueOf(b)
	case reflect.Int:
		s, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		return reflect.ValueOf(s)
	case reflect.String:
		return reflect.ValueOf(v)
	default:
		panic(fmt.Errorf("type %s is not support", typ.Kind().String()))
	}
}
