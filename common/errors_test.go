// Copyright(c) 2018 DSiSc Group All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	fmt "fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorPanic(t *testing.T) {
	type pnk struct {
		msg string
	}

	capturePanic := func() (err Error) {
		defer func() {
			if r := recover(); r != nil {
				err = ErrorWrap(r, "This is the message in ErrorWrap(r, message).")
			}
		}()
		panic(pnk{"something"})
	}

	var err = capturePanic()

	assert.Equal(t, pnk{"something"}, err.Data())
	assert.Equal(t, "Error{{something}}", fmt.Sprintf("%v", err))
	assert.Contains(t, fmt.Sprintf("%#v", err), "This is the message in ErrorWrap(r, message).")
	assert.Contains(t, fmt.Sprintf("%#v", err), "Stack Trace:\n    0")
}

func TestErrorWrapSomething(t *testing.T) {

	var err = ErrorWrap("something", "formatter%v%v", 0, 1)

	assert.Equal(t, "something", err.Data())
	assert.Equal(t, "Error{something}", fmt.Sprintf("%v", err))
	assert.Regexp(t, `formatter01\n`, fmt.Sprintf("%#v", err))
	assert.Contains(t, fmt.Sprintf("%#v", err), "Stack Trace:\n    0")
}

func TestErrorWrapNothing(t *testing.T) {

	var err = ErrorWrap(nil, "formatter%v%v", 0, 1)

	assert.Equal(t,
		FmtError{"formatter%v%v", []interface{}{0, 1}},
		err.Data())
	assert.Equal(t, "Error{formatter01}", fmt.Sprintf("%v", err))
	assert.Contains(t, fmt.Sprintf("%#v", err), `Data: common.FmtError{format:"formatter%v%v", args:[]interface {}{0, 1}}`)
	assert.Contains(t, fmt.Sprintf("%#v", err), "Stack Trace:\n    0")
}

func TestErrorNewError(t *testing.T) {

	var err = NewError("formatter%v%v", 0, 1)

	assert.Equal(t,
		FmtError{"formatter%v%v", []interface{}{0, 1}},
		err.Data())
	assert.Equal(t, "Error{formatter01}", fmt.Sprintf("%v", err))
	assert.Contains(t, fmt.Sprintf("%#v", err), `Data: common.FmtError{format:"formatter%v%v", args:[]interface {}{0, 1}}`)
	assert.NotContains(t, fmt.Sprintf("%#v", err), "Stack Trace")
}

func TestErrorNewErrorWithStacktrace(t *testing.T) {

	var err = NewError("formatter%v%v", 0, 1).Stacktrace()

	assert.Equal(t,
		FmtError{"formatter%v%v", []interface{}{0, 1}},
		err.Data())
	assert.Equal(t, "Error{formatter01}", fmt.Sprintf("%v", err))
	assert.Contains(t, fmt.Sprintf("%#v", err), `Data: common.FmtError{format:"formatter%v%v", args:[]interface {}{0, 1}}`)
	assert.Contains(t, fmt.Sprintf("%#v", err), "Stack Trace:\n    0")
}

func TestErrorNewErrorWithTrace(t *testing.T) {

	var err = NewError("formatter%v%v", 0, 1)
	err.Trace(0, "trace %v", 1)
	err.Trace(0, "trace %v", 2)
	err.Trace(0, "trace %v", 3)

	assert.Equal(t,
		FmtError{"formatter%v%v", []interface{}{0, 1}},
		err.Data())
	assert.Equal(t, "Error{formatter01}", fmt.Sprintf("%v", err))
	assert.Contains(t, fmt.Sprintf("%#v", err), `Data: common.FmtError{format:"formatter%v%v", args:[]interface {}{0, 1}}`)
	dump := fmt.Sprintf("%#v", err)
	assert.NotContains(t, dump, "Stack Trace")
	assert.Regexp(t, `common/errors_test\.go:[0-9]+ - trace 1`, dump)
	assert.Regexp(t, `common/errors_test\.go:[0-9]+ - trace 2`, dump)
	assert.Regexp(t, `common/errors_test\.go:[0-9]+ - trace 3`, dump)
}

func TestErrorWrapError(t *testing.T) {
	var err1 error = NewError("my message")
	var err2 error = ErrorWrap(err1, "another message")
	assert.Equal(t, err1, err2)
}
