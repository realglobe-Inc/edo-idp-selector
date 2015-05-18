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

package idp

import (
	"fmt"
	"reflect"
	"testing"
)

var (
	test_elem  = newElement(test_id, test_names, test_authUri, test_tokUri, test_acntUri, test_coopFrUri, test_keys)
	test_elem2 = newElement(
		"https://idp2.example.org",
		map[string]string{"": "Genuine ID Provider", "ja": "真の ID プロバイダ"},
		"https://idp2.example.org/auth",
		"https://idp2.example.org/token",
		"https://idp2.example.org/userinfo",
		"https://idp2.example.org/cooperation/from",
		nil,
	)
)

// test_elem と test_elem2 が保存されていることが前提。
func testDb(t *testing.T, db Db) {
	if elem, err := db.Get(test_elem.Id() + "a"); err != nil {
		t.Fatal(err)
	} else if elem != nil {
		t.Fatal(elem)
	} else if elem, err := db.Get(test_elem.Id()); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(elem, test_elem) {
		t.Error(elem)
		t.Fatal(test_elem)
	} else if elems, err := db.Search(map[string]string{"issuer": `aaaaa`}); err != nil {
		t.Fatal(err)
	} else if len(elems) > 0 {
		t.Fatal(elems)
	} else if elems, err := db.Search(map[string]string{"issuer": `https://idp\.example\.org`}); err != nil {
		t.Fatal(err)
	} else if len(elems) != 1 {
		t.Fatal(elems)
	} else if !reflect.DeepEqual(elems[0], test_elem) {
		t.Error(elems[0])
		t.Fatal(test_elem)
	} else if elems, err := db.Search(map[string]string{"issuer_name#ja": "真の"}); err != nil {
		t.Fatal(err)
	} else if len(elems) != 1 {
		t.Fatal(elems)
	} else if !reflect.DeepEqual(elems[0], test_elem2) {
		t.Error(fmt.Sprintf("%#v", elems[0]))
		t.Fatal(fmt.Sprintf("%#v", test_elem2))
	} else if elems, err := db.Search(map[string]string{"authorization_endpoint": `idp\.`}); err != nil {
		t.Fatal(err)
	} else if len(elems) != 1 {
		t.Fatal(elems)
	} else if elems, err := db.Search(map[string]string{"token_endpoint": `idp\.`}); err != nil {
		t.Fatal(err)
	} else if len(elems) != 1 {
		t.Fatal(elems)
	} else if elems, err := db.Search(map[string]string{"userinfo_endpoint": `idp\.`}); err != nil {
		t.Fatal(err)
	} else if len(elems) != 1 {
		t.Fatal(elems)
	} else if !reflect.DeepEqual(elems[0], test_elem) {
		t.Error(elems[0])
		t.Fatal(test_elem)
	} else if elems, err := db.Search(map[string]string{"cooperation_from_endpoint": `https://idp[0-9]*\.example\.org/cooperation/from`}); err != nil {
		t.Fatal(err)
	} else if len(elems) != 2 {
		t.Fatal(elems)
	} else if (!reflect.DeepEqual(elems[0], test_elem) && !reflect.DeepEqual(elems[1], test_elem)) ||
		!reflect.DeepEqual(elems[0], test_elem2) && !reflect.DeepEqual(elems[1], test_elem2) {
		t.Error(elems)
		t.Error(test_elem)
		t.Fatal(test_elem2)
	}
}
