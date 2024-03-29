package goenv

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var Environ = os.Environ

func analyse(target interface{}) ([]pointer, error) {
	typ := reflect.TypeOf(target)
	if typ.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("expect pointer to struct")
	}
	typ = typ.Elem()
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expect pointer to struct")
	}

	var curr pointer
	analyzer := analyzer{}
	analyzer.analyzeRec(typ, curr)

	return analyzer.output, nil
}

type analyzer struct {
	output []pointer
}

func (a *analyzer) analyzeRec(typ reflect.Type, curr pointer) {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	switch typ.Kind() {
	case reflect.Bool:
		fallthrough
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		fallthrough
	case reflect.String:
		a.output = append(a.output, curr)
	case reflect.Struct:
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			hint := field.Tag.Get("goenv")
			if hint == "-" {
				continue
			}
			var next pointer
			// TODO: Validate hint
			if hint != "" {
				next = curr.StructField(strings.ToUpper(hint), i)
			} else {
				next = curr.StructField(strings.ToUpper(field.Name), i)
			}
			a.analyzeRec(field.Type, next)
		}
	case reflect.Slice:
		typ = typ.Elem()
		next := curr.SliceIndex(0)
		a.analyzeRec(typ, next)
	}
}

func specialize(ptrs []pointer, env environment) []pointer {
	var output []pointer
	for _, ptr := range ptrs {
		for name := range env {
			p, ok := ptr.specialize(name)
			if ok {
				output = append(output, p)
			}
		}
	}
	sortPointers(output)
	return output
}

func assign(ptrs []pointer, env environment, target interface{}) error {
	for _, ptr := range ptrs {
		value := reflect.ValueOf(target)
		for _, part := range ptr {
			for value.Kind() == reflect.Ptr {
				if value.IsNil() {
					value.Set(reflect.New(value.Type().Elem()))
				}
				value = value.Elem()
			}
			if part.Name == "" {
				if part.SliceIndex >= value.Len() {
					newValue := reflect.MakeSlice(value.Type(), part.SliceIndex+1, part.SliceIndex+1)
					reflect.Copy(newValue, value)
					value.Set(newValue)
				}
				value = value.Index(part.SliceIndex)
			} else {
				value = value.Field(part.FieldIndex)
			}
		}
		for value.Kind() == reflect.Ptr {
			if value.IsNil() {
				value.Set(reflect.New(value.Type().Elem()))
			}
			value = value.Elem()
		}
		envVarName := ptr.String()
		envVarValue, ok := env.LookupEnv(envVarName)
		if !ok {
			return fmt.Errorf("unexpected lookup: %v", envVarName)
		}
		switch value.Kind() {
		case reflect.Bool:
			b, err := strconv.ParseBool(envVarValue)
			if err != nil {
				return err
			}
			value.SetBool(b)
		case reflect.Int:
			i, err := strconv.ParseInt(envVarValue, 10, 0)
			if err != nil {
				return err
			}
			value.SetInt(i)
		case reflect.Int8:
			i, err := strconv.ParseInt(envVarValue, 10, 8)
			if err != nil {
				return err
			}
			value.SetInt(i)
		case reflect.Int16:
			i, err := strconv.ParseInt(envVarValue, 10, 16)
			if err != nil {
				return err
			}
			value.SetInt(i)
		case reflect.Int32:
			i, err := strconv.ParseInt(envVarValue, 10, 32)
			if err != nil {
				return err
			}
			value.SetInt(i)
		case reflect.Int64:
			i, err := strconv.ParseInt(envVarValue, 10, 64)
			if err != nil {
				return err
			}
			value.SetInt(i)
		case reflect.Uint:
			u, err := strconv.ParseUint(envVarValue, 10, 0)
			if err != nil {
				return err
			}
			value.SetUint(u)
		case reflect.Uint8:
			u, err := strconv.ParseUint(envVarValue, 10, 8)
			if err != nil {
				return err
			}
			value.SetUint(u)
		case reflect.Uint16:
			u, err := strconv.ParseUint(envVarValue, 10, 16)
			if err != nil {
				return err
			}
			value.SetUint(u)
		case reflect.Uint32:
			u, err := strconv.ParseUint(envVarValue, 10, 32)
			if err != nil {
				return err
			}
			value.SetUint(u)
		case reflect.Uint64:
			u, err := strconv.ParseUint(envVarValue, 10, 64)
			if err != nil {
				return err
			}
			value.SetUint(u)
		case reflect.Float32:
			f, err := strconv.ParseFloat(envVarValue, 32)
			if err != nil {
				return err
			}
			value.SetFloat(f)
		case reflect.Float64:
			f, err := strconv.ParseFloat(envVarValue, 64)
			if err != nil {
				return err
			}
			value.SetFloat(f)
		case reflect.String:
			value.SetString(envVarValue)
		}
	}
	return nil
}

// TODO: Support default? (For struct field and slice element)
func Populate(prefix string, target interface{}) error {
	ptrs, err := analyse(target)
	if err != nil {
		return err
	}
	env := parseEnvironment(prefix, Environ)
	ptrs = specialize(ptrs, env)
	err = assign(ptrs, env, target)
	return err
}
