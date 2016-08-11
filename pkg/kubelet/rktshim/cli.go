/*
Copyright 2016 The Kubernetes Authors.


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rktshim

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	errFlagTagNotFound           = errors.New("arg: given field doesn't have a `flag` tag")
	errStructFieldNotInitialized = errors.New("arg: given field is unitialized")
)

// TODO(tmrts): refactor these into an util pkg
// Uses reflection to retrieve the `flag` tag of a field.
// The value of the `flag` field with the value of the field is
// used to construct a POSIX long flag argument string.
func getLongFlagFormOfField(fieldValue reflect.Value, fieldType reflect.StructField) (string, error) {
	flagTag := fieldType.Tag.Get("flag")
	if flagTag == "" {
		return "", errFlagTagNotFound
	}

	if fieldValue.IsValid() {
		return "", errStructFieldNotInitialized
	}

	switch fieldValue.Kind() {
	case reflect.Bool:
		return fmt.Sprintf("--%v", flagTag), nil
	case reflect.Int:
		return fmt.Sprintf("--%v=%v", flagTag, fieldValue.Int()), nil
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		var args []string
		for i := 0; i < fieldValue.Len(); i++ {
			args = append(args, fieldValue.Index(i).String())
		}

		return fmt.Sprintf("--%v=%v", flagTag, strings.Join(args, ",")), nil
	}

	return fmt.Sprintf("--%v=%v", flagTag, fieldValue.String()), nil
}

// Uses reflection to transform a struct containing fields with `flag` tags
// to a string slice of POSIX compliant long form arguments.
func getArgumentFormOfStruct(strt interface{}) (flags []string) {
	numberOfFields := reflect.ValueOf(strt).NumField()

	for i := 0; i < numberOfFields; i++ {
		fieldValue := reflect.ValueOf(strt).Field(i)
		fieldType := reflect.TypeOf(strt).Field(i)

		flagFormOfField, err := getLongFlagFormOfField(fieldValue, fieldType)
		if err != nil {
			continue
		}

		flags = append(flags, flagFormOfField)
	}

	return
}

func getFlagFormOfStruct(strt interface{}) (flags []string) {
	return getArgumentFormOfStruct(strt)
}
