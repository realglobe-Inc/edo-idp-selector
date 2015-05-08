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

package web

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestDirectDb(t *testing.T) {
	path := "/a/b/c"
	mux := http.NewServeMux()
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Write(test_data)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	// URI が test_elem と違うので testDb は使えない。
	elem := newElement(server.URL+path, test_data)

	db := NewDirectDb()
	if el, err := db.Get(elem.Uri() + "a"); err != nil {
		t.Fatal(err)
	} else if el != nil {
		t.Fatal(el)
	} else if el, err := db.Get(elem.Uri()); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(el, elem) {
		t.Error(el)
		t.Fatal(elem)
	}
}
