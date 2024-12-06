package utils

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

func ParseForm(toFill interface{}, request *http.Request) error {
	v := reflect.ValueOf(toFill)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct")
	}
	v = v.Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("form")
		if tag == "" {
			continue
		}

		fieldValue := v.Field(i)
		if !fieldValue.CanSet() {
			continue
		}

		switch fieldValue.Kind() {
		case reflect.String:
			fieldValue.SetString(request.FormValue(tag))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value, err := strconv.ParseInt(request.FormValue(tag), 10, 64)
			if err != nil {
				return err
			}
			fieldValue.SetInt(value)
		case reflect.Float32, reflect.Float64:
			value, err := strconv.ParseFloat(request.FormValue(tag), 64)
			if err != nil {
				return err
			}
			fieldValue.SetFloat(value)
		case reflect.Pointer:
			value := request.FormValue(tag)

			if value == "" {
				continue
			}

			fieldType := fieldValue.Type().Elem()
			newValue := reflect.New(fieldType).Elem()

			switch fieldType.Kind() {
			case reflect.String:
				newValue.SetString(value)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				intValue, err := strconv.ParseInt(request.FormValue(tag), 10, 64)
				if err != nil {
					return err
				}
				newValue.SetInt(intValue)
			case reflect.Float32, reflect.Float64:
				value, err := strconv.ParseFloat(request.FormValue(tag), 64)
				if err != nil {
					return err
				}
				newValue.SetFloat(value)
			}
			fieldValue.Set(newValue.Addr())
		}
	}
	return nil
}
