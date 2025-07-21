package python

import "C"
import (
	"fmt"
	"reflect"

	cpy3 "go.nhat.io/cpy/v3"
)

// Marshaler is the interface implemented by types that can marshal themselves into a Python object.
type Marshaler interface {
	MarshalPyObject() *Object
}

// Unmarshaler is the interface implemented by types that can unmarshal a Python object of themselves.
type Unmarshaler interface {
	UnmarshalPyObject(o *Object) error
}

// Marshal returns the Python object for v.
func Marshal(v any) (*Object, error) { //nolint: cyclop,funlen,gocyclo
	if v, ok := v.(Marshaler); ok {
		return v.MarshalPyObject(), nil
	}

	if v := reflect.ValueOf(v); v.Kind() == reflect.Ptr && v.IsNil() {
		return nil, nil //nolint: nilnil
	}

	switch v := v.(type) {
	case *cpy3.PyObject:
		return NewObject(v), nil

	case *Object:
		return v, nil

	case Objector:
		return v.AsObject(), nil

	case PyObjector:
		return NewObject(v.PyObject()), nil

	case bool:
		return NewBool(v), nil

	case string:
		return NewString(v), nil

	case int:
		return NewInt(v), nil

	case uint:
		return NewUint(v), nil

	case int8:
		return NewInt(int(v)), nil

	case uint8:
		return NewUint(uint(v)), nil

	case int16:
		return NewInt(int(v)), nil

	case uint16:
		return NewUint(uint(v)), nil

	case int32:
		return NewInt(int(v)), nil

	case uint32:
		return NewUint(uint(v)), nil

	case int64:
		return NewInt64(v), nil

	case uint64:
		return NewUint64(v), nil

	case float32:
		return NewFloat64(float64(v)), nil

	case float64:
		return NewFloat64(v), nil
	}

	rv := reflect.Indirect(reflect.ValueOf(v))
	if rv.Kind() == reflect.Slice {
		return marshalSlice(rv), nil
	}

	return nil, fmt.Errorf("cannot marshal value of %T to python object", v) //nolint: err113
}

// MustMarshal returns the Python object for v or panics if an error occurs.
func MustMarshal(v any) *Object {
	o, err := Marshal(v)
	if err != nil {
		panic(err)
	}

	return o
}

func marshalSlice(v reflect.Value) *Object {
	l := make([]any, v.Cap(), v.Len())

	for i := range v.Len() {
		l[i] = v.Index(i).Interface()
	}

	return NewListFromAny(l...).AsObject()
}

// An InvalidUnmarshalError describes an invalid argument passed to [Unmarshal].
// (The argument to [Unmarshal] must be a non-nil pointer).
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "python3: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Pointer {
		return fmt.Sprintf("python3: Unmarshal(non-pointer %s)", e.Type.String())
	}

	return fmt.Sprintf("python3: Unmarshal(nil %s)", e.Type.String())
}

// An UnmarshalTypeError describes a python Object that was not appropriate for a value of a specific Go type.
type UnmarshalTypeError struct {
	Value  string       // Description of the value - "bool", "array", "number -5".
	Type   reflect.Type // Type of Go value it could not be assigned to.
	Struct string       // Name of the struct type containing the field.
	Field  string       // The full path from root node to the field.
}

func (e *UnmarshalTypeError) Error() string {
	if e.Struct != "" || e.Field != "" {
		return "python3: cannot unmarshal " + e.Value + " into Go struct field " + e.Struct + "." + e.Field + " of type " + e.Type.String()
	}

	return "python3: cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()
}

// Unmarshal converts the Python object to a value of the same type as v.
func Unmarshal(o *Object, v any) error { //nolint: cyclop,funlen,gocognit,gocyclo
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	if u, ok := v.(Unmarshaler); ok {
		return u.UnmarshalPyObject(o)
	}

	irv := reflect.Indirect(rv)
	kind := irv.Kind()
	targetKind := kind

	if irv.Type() == reflect.TypeOf((*any)(nil)).Elem() {
		targetKind = objectKind(o)
	}

	switch targetKind {
	case reflect.String:
		v, err := unmarshalString(o)
		if err != nil {
			return err
		}

		if kind != targetKind {
			irv.Set(reflect.ValueOf(v))
		} else {
			irv.SetString(v)
		}

		return nil

	case reflect.Bool:
		v, err := unmarshalBool(o)
		if err != nil {
			return err
		}

		if kind != targetKind {
			irv.Set(reflect.ValueOf(v))
		} else {
			irv.SetBool(v)
		}

		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := unmarshalInt(o)
		if err != nil {
			return err
		}

		if kind != targetKind {
			irv.Set(reflect.ValueOf(v))
		} else {
			irv.SetInt(v)
		}

		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := unmarshalUint(o)
		if err != nil {
			return err
		}

		if kind != targetKind {
			irv.Set(reflect.ValueOf(v))
		} else {
			irv.SetUint(v)
		}

		return nil

	case reflect.Float32, reflect.Float64:
		v, err := unmarshalFloat64(o)
		if err != nil {
			return err
		}

		if kind != targetKind {
			irv.Set(reflect.ValueOf(v))
		} else {
			irv.SetFloat(v)
		}

		return nil

	case reflect.Slice:
		return unmarshalSlice(o, irv)

	case reflect.Pointer:
		p := reflect.New(irv.Type().Elem())

		if err := Unmarshal(o, p.Interface()); err != nil {
			return err
		}

		irv.Set(p)

		return nil

	default:
	}

	return &UnmarshalTypeError{Value: TypeName(o), Type: irv.Type()}
}

// UnmarshalAs converts the Python object to a value of the same type as T.
func UnmarshalAs[T any](o *Object) (T, error) {
	var v T

	if err := Unmarshal(o, &v); err != nil {
		return v, err
	}

	return v, nil
}

// MustUnmarshal converts the Python object to a value of the same type as v or panics if an error occurs.
func MustUnmarshal(o *Object, v any) {
	if err := Unmarshal(o, v); err != nil {
		panic(err)
	}
}

// MustUnmarshalAs converts the Python object to a value of the same type as T or panics if an error occurs.
func MustUnmarshalAs[T any](o *Object) T {
	var v T

	MustUnmarshal(o, &v)

	return v
}

func unmarshalString(o *Object) (string, error) {
	if !IsString(o) {
		return "", &UnmarshalTypeError{Value: TypeName(o), Type: reflect.TypeFor[string]()}
	}

	return o.String(), nil
}

func unmarshalBool(o *Object) (bool, error) {
	if !IsBool(o) {
		return false, &UnmarshalTypeError{Value: TypeName(o), Type: reflect.TypeFor[bool]()}
	}

	return AsBool(o), nil
}

func unmarshalInt(o *Object) (int64, error) {
	if !IsInt(o) {
		return 0, &UnmarshalTypeError{Value: TypeName(o), Type: reflect.TypeFor[int64]()}
	}

	return AsInt64(o), nil
}

func unmarshalUint(o *Object) (uint64, error) {
	if !IsInt(o) {
		return 0, &UnmarshalTypeError{Value: TypeName(o), Type: reflect.TypeFor[uint64]()}
	}

	return AsUint64(o), nil
}

func unmarshalFloat64(o *Object) (float64, error) {
	if !IsFloat(o) {
		return 0, &UnmarshalTypeError{Value: TypeName(o), Type: reflect.TypeFor[float64]()}
	}

	return AsFloat64(o), nil
}

func unmarshalSlice(o *Object, dest reflect.Value) error {
	if !IsList(o) && !IsTuple(o) {
		return &UnmarshalTypeError{Value: TypeName(o), Type: dest.Type()}
	}

	v := reflect.MakeSlice(dest.Type(), o.Length(), o.Length())
	defers := make([]func(), 0, o.Length())

	defer func() {
		for _, def := range defers {
			def()
		}
	}()

	for i := range o.Length() {
		item := o.GetItem(i)

		defers = append(defers, item.DecRef)

		if err := Unmarshal(item, v.Index(i).Addr().Interface()); err != nil {
			return err
		}
	}

	dest.Set(v)

	return nil
}

func objectKind(o *Object) reflect.Kind {
	if IsBool(o) {
		return reflect.Bool
	}

	if IsInt(o) {
		return reflect.Int
	}

	if IsFloat(o) {
		return reflect.Float64
	}

	if IsString(o) {
		return reflect.String
	}

	if IsList(o) || IsTuple(o) {
		return reflect.Slice
	}

	return reflect.Invalid
}
