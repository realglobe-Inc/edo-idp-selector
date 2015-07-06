// Copyright 2015 realglobe, Inc.
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

package error

import (
	"net/http"

	"github.com/realglobe-Inc/go-lib/erro"
)

// OAuth 2.0 や OpenID Connect 1.0 でのエラーレスポンスに必要な情報を含むエラー。
type Error struct {
	// error の値。
	errCod string
	// error_description の値。
	errDesc string

	// 直接返すときの HTTP ステータス。
	stat int

	// 原因になったエラー。
	cause error
}

// stat が 0 の場合、代わりに http.StatusInternalServerError が入る。
func New(errCod string, errDesc string, stat int, cause error) *Error {
	if stat <= 0 {
		stat = http.StatusInternalServerError
	}
	return &Error{
		errCod:  errCod,
		errDesc: errDesc,
		stat:    stat,
		cause:   cause,
	}
}

func (this *Error) Error() string {
	pref := ""
	if this.cause != nil {
		pref += this.cause.Error() + "\ncaused "
	}
	return pref + this.errCod + ": " + this.errDesc
}

func (this *Error) ErrorCode() string {
	return this.errCod
}

func (this *Error) ErrorDescription() string {
	return this.errDesc
}

func (this *Error) Status() int {
	return this.stat
}

func (this *Error) Cause() error {
	return this.cause
}

// 通常のエラーから変換する。
func From(err error) *Error {
	if e, ok := err.(*Error); ok {
		return e
	}
	e2 := erro.Unwrap(err)
	if e, ok := e2.(*Error); ok {
		return New(e.ErrorCode(), e.ErrorDescription(), e.Status(), err)
	} else {
		return New(Server_error, e2.Error(), http.StatusInternalServerError, err)
	}
}
