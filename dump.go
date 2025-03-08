// Copyright 2014 The godump Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package godump

import (
	"fmt"
	"reflect"
	"strconv"
)

type variable struct {
	// Output dump string
	Out string

	// Indent counter
	indent int64
}

// Stringer is implemented by any value that has a String method
type Stringer interface {
	String() string
}

// dump outputs a formatted string for the given value and name.
// The string is indented with the current indent level, which is
// incremented before the dump and decremented after.
//
// The output format is as follows:
//
//   - For arrays and slices, the type and length are printed,
//     followed by the elements.
//   - For maps, the type and key type are printed, followed by the
//     key-value pairs.
//   - For pointers, the type is printed, followed by the result of
//     calling dump on the value.
//   - For structs, the type is printed, followed by the fields.
//   - For all other types, the value is printed using the fmt
//     package's formatting rules.
func (v *variable) dump(val reflect.Value, name string) {
	v.indent++

	if val.IsValid() && val.CanInterface() {
		typ := val.Type()

		switch typ.Kind() {
		case reflect.Array, reflect.Slice:
			v.printType(name, val.Interface())
			l := val.Len()
			for i := 0; i < l; i++ {
				v.dump(val.Index(i), strconv.Itoa(i))
			}
		case reflect.Map:
			v.printType(name, val.Interface())
			//l := val.Len()
			keys := val.MapKeys()
			for _, k := range keys {
				v.dump(val.MapIndex(k), k.Interface().(string))
			}
		case reflect.Ptr:
			v.printType(name, val.Interface())
			v.dump(val.Elem(), name)
		case reflect.Struct:

			v.printType(name, val.Interface())
			for i := 0; i < typ.NumField(); i++ {
				field := typ.Field(i)
				v.dump(val.FieldByIndex([]int{i}), field.Name)
			}
		default:
			v.printValue(name, val.Interface())
		}
	} else {
		v.printValue(name, "")
	}

	v.indent--
}

// printType writes a type information of value to the dump string.
//
// It adds an indentation prefix, name, type of value, and a newline character
// to the dump string. If value implements Stringer interface, it also adds a
// string representation of value to the dump string.
func (v *variable) printType(name string, vv interface{}) {
	v.printIndent()
	_, ok := vv.(Stringer)
	if ok {
		v.Out = fmt.Sprintf("%s%s(%T) %s\n", v.Out, name, vv, vv)
		return
	}

	v.Out = fmt.Sprintf("%s%s(%T)\n", v.Out, name, vv)
}

// printValue writes a value information of value to the dump string.
//
// It adds an indentation prefix, name, type of value, value formatted with
// %#v, and a newline character to the dump string.
func (v *variable) printValue(name string, vv interface{}) {
	v.printIndent()
	v.Out = fmt.Sprintf("%s%s(%T) %#v\n", v.Out, name, vv, vv)
}

// printIndent adds an indentation prefix to the dump string.
//
// The number of space characters in the prefix is twice the value of
// v.indent.
func (v *variable) printIndent() {
	var i int64
	for i = 0; i < v.indent; i++ {
		v.Out = fmt.Sprintf("%s  ", v.Out)
	}
}

// Print to standard out the value that is passed as the argument with indentation.
// Pointers are dereferenced.
func Dump(v interface{}) {
	val := reflect.ValueOf(v)
	dump := &variable{indent: -1}
	dump.dump(val, "")
	fmt.Print(dump.Out)
}

// Return the value that is passed as the argument with indentation.
// Pointers are dereferenced.
func Sdump(v interface{}) string {
	val := reflect.ValueOf(v)
	dump := &variable{indent: -1}
	dump.dump(val, "")
	return dump.Out
}
