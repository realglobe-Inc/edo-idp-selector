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
	"errors"
	"net/http"
	"testing"

	"github.com/realglobe-Inc/go-lib/erro"
)

const (
	test_errCod  = Access_denied
	test_errDesc = "you are denied by everyone"
	test_stat    = http.StatusForbidden
	test_msg     = "the end"
)

var (
	test_cause = erro.New(test_msg)
)

func TestError(t *testing.T) {

	if err := New(test_errCod, test_errDesc, test_stat, test_cause); err.ErrorCode() != test_errCod {
		t.Error(err.ErrorCode())
		t.Fatal(test_errCod)
	} else if err.ErrorDescription() != test_errDesc {
		t.Error(err.ErrorDescription())
		t.Fatal(test_errDesc)
	} else if err.Status() != test_stat {
		t.Error(err.Status())
		t.Fatal(test_stat)
	} else if err.Cause() != test_cause {
		t.Error(err.Cause())
		t.Fatal(test_cause)
	}
}

func TestErrorFrom(t *testing.T) {
	err := New(test_errCod, test_errDesc, test_stat, test_cause)
	if err2 := From(err); err2 != err {
		t.Error(err2)
		t.Fatal(err)
	} else {
		//t.Error(err)
	}

	if err := From(errors.New(test_msg)); err.ErrorCode() != Server_error {
		t.Error(err.ErrorCode())
		t.Fatal(Server_error)
	} else if err.ErrorDescription() != test_msg {
		t.Error(err.ErrorDescription())
		t.Fatal(test_msg)
	} else if err.Status() != http.StatusInternalServerError {
		t.Error(err.Status())
		t.Fatal(http.StatusInternalServerError)
	} else if err.Cause() == nil {
		t.Fatal("no cause")
	} else {
		//t.Error(err)
	}

	if err := From(erro.Wrap(New(test_errCod, test_errDesc, test_stat, test_cause))); err.ErrorCode() != test_errCod {
		t.Error(err.ErrorCode())
		t.Fatal(test_errCod)
	} else if err.ErrorDescription() != test_errDesc {
		t.Error(err.ErrorDescription())
		t.Fatal(test_errDesc)
	} else if err.Status() != test_stat {
		t.Error(err.Status())
		t.Fatal(test_stat)
	} else if err.Cause() == test_cause {
		t.Fatal("same cause")
	} else {
		//t.Error(err)
	}
}
