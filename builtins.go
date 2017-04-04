// Package validator implements value validations
//
// Copyright 2014 Roberto Teixeira <robteix@robteix.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validator

import (
	"reflect"
	"regexp"
	"strconv"
	"unicode/utf8"
)

// nonzero tests whether a variable value non-zero
// as defined by the golang spec.
func nonzero(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	switch st.Kind() {
	case reflect.String:
		if utf8.RuneCount([]byte(st.String())) == 0 {
			return ErrZeroValueEmpty
		}
	case reflect.Ptr, reflect.Interface:
		if st.IsNil() {
			return ErrZeroValueEmpty
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		if st.Len() == 0 {
			return ErrZeroValueEmpty
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if st.Int() == 0 {
			return ErrZeroValueNumber
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if st.Uint() == 0 {
			return ErrZeroValueNumber
		}
	case reflect.Float32, reflect.Float64:
		if st.Float() == 0 {
			return ErrZeroValueNumber
		}
	case reflect.Bool:
		if !st.Bool() {
			return ErrZeroValueBool
		}
	case reflect.Invalid:
		return ErrZeroValue
	case reflect.Struct:
		return nil
	default:
		return ErrUnsupported
	}

	return nil
}

// length tests whether a variable's length is equal to a given
// value. For strings it tests the number of characters whereas
// for maps and slices it tests the number of items.
func length(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	switch st.Kind() {
	case reflect.String:
		p, err := asInt(param)
		if err != nil {
			return ErrBadParameter
		}
		actual := utf8.RuneCount([]byte(st.String()))
		if int64(actual) != p {
			return ErrLenString(p, actual)
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		p, err := asInt(param)
		if err != nil {
			return ErrBadParameter
		}
		actual := st.Len()
		if int64(actual) != p {
			return ErrLenArray(p, actual)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := asInt(param)
		if err != nil {
			return ErrBadParameter
		}
		actual := st.Int()
		if actual != p {
			return ErrLenInt(p, actual)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := asUint(param)
		if err != nil {
			return ErrBadParameter
		}
		actual := st.Uint()
		if actual != p {
			return ErrLenInt(int64(p), int64(actual))
		}
	case reflect.Float32, reflect.Float64:
		p, err := asFloat(param)
		if err != nil {
			return ErrBadParameter
		}
		actual := st.Float()
		if actual != p {
			return ErrLenFloat(p, actual)
		}
	case reflect.Ptr:
		return nil
	default:
		return ErrUnsupported
	}

	return nil
}

// min tests whether a variable value is larger or equal to a given
// number. For number types, it's a simple lesser-than test; for
// strings it tests the number of characters whereas for maps
// and slices it tests the number of items.
func min(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	switch st.Kind() {
	case reflect.String:
		p, err := asInt(param)
		if err != nil {
			return ErrBadParameter
		}
		actual := utf8.RuneCount([]byte(st.String()))
		if int64(actual) < p {
			return ErrMinString(p, actual)
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		p, err := asInt(param)
		if err != nil {
			return ErrBadParameter
		}
		actual := st.Len()
		if int64(actual) < p {
			return ErrMinArray(p, actual)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := asInt(param)
		if err != nil {
			return ErrBadParameter
		}
		actual := st.Int()
		if actual < p {
			return ErrMinInt(int64(p), int64(actual))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := asUint(param)
		if err != nil {
			return ErrBadParameter
		}
		actual := st.Uint()
		if actual < p {
			return ErrMinInt(int64(p), int64(actual))
		}
	case reflect.Float32, reflect.Float64:
		p, err := asFloat(param)
		if err != nil {
			return ErrBadParameter
		}
		actual := st.Float()
		if actual < p {
			return ErrMinFloat(p, actual)
		}
	case reflect.Ptr:
		return nil
	default:
		return ErrUnsupported
	}

	return nil
}

// max tests whether a variable value is lesser than a given
// value. For numbers, it's a simple lesser-than test; for
// strings it tests the number of characters whereas for maps
// and slices it tests the number of items.
func max(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	switch st.Kind() {
	case reflect.String:
		p, err := asInt(param)
		if err != nil {
			return ErrBadParameter
		}
		actual := utf8.RuneCount([]byte(st.String()))
		if int64(actual) > p {
			return ErrMaxString(p, actual)
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		p, err := asInt(param)
		if err != nil {
			return ErrBadParameter
		}
		actual := st.Len()
		if int64(actual) > p {
			return ErrMaxArray(p, actual)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := asInt(param)
		if err != nil {
			return ErrBadParameter
		}
		actual := st.Int()
		if actual > p {
			return ErrMaxInt(p, actual)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := asUint(param)
		if err != nil {
			return ErrBadParameter
		}
		actual := st.Uint()
		if actual > p {
			return ErrMaxInt(int64(p), int64(actual))
		}
	case reflect.Float32, reflect.Float64:
		p, err := asFloat(param)
		if err != nil {
			return ErrBadParameter
		}
		actual := st.Float()
		if actual > p {
			return ErrMaxFloat(p, actual)
		}
	case reflect.Ptr:
		return nil
	default:
		return ErrUnsupported
	}
	return nil
}

// regex is the builtin validation function that checks
// whether the string variable matches a regular expression
func regex(v interface{}, param string) error {
	s, ok := v.(string)
	if !ok {
		return ErrUnsupported
	}

	re, err := regexp.Compile(param)
	if err != nil {
		return ErrBadParameter
	}

	if !re.MatchString(s) {
		return ErrRegexpDetailed(param)
	}
	return nil
}

// asInt retuns the parameter as a int64
// or panics if it can't convert
func asInt(param string) (int64, error) {
	i, err := strconv.ParseInt(param, 0, 64)
	if err != nil {
		return 0, ErrBadParameter
	}
	return i, nil
}

// asUint retuns the parameter as a uint64
// or panics if it can't convert
func asUint(param string) (uint64, error) {
	i, err := strconv.ParseUint(param, 0, 64)
	if err != nil {
		return 0, ErrBadParameter
	}
	return i, nil
}

// asFloat retuns the parameter as a float64
// or panics if it can't convert
func asFloat(param string) (float64, error) {
	i, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return 0.0, ErrBadParameter
	}
	return i, nil
}
