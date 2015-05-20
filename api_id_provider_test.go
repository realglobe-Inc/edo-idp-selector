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

package main

import (
	"encoding/json"
	idpdb "github.com/realglobe-Inc/edo-idp-selector/database/idp"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIdProviderApi(t *testing.T) {
	sys := newTestSystem(nil, []idpdb.Element{
		test_idp1,
		test_idp2,
		test_idp3,
	}, nil)

	r, err := http.NewRequest("GET", "https://selector.example.org/api/info/issuer", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	if err := sys.idProviderApi(w, r); err != nil {
		t.Fatal(err)
	} else if w.Code != http.StatusOK {
		t.Error(w.Code)
		t.Fatal(http.StatusOK)
	} else if w.HeaderMap.Get("Content-Type") != "application/json" {
		t.Error(w.HeaderMap.Get("Content-Type"))
		t.Fatal("application/json")
	}

	data, _ := ioutil.ReadAll(w.Body)
	var buff []struct {
		Issuer        string
		Issuer_name   string
		Issuer_nameJa string `json:"issuer_name#ja"`
	}
	if err := json.Unmarshal(data, &buff); err != nil {
		t.Fatal(err)
	} else if len(buff) != 3 {
		t.Fatal(buff)
	}
	m := map[string]idpdb.Element{
		test_idp1.Id(): test_idp1,
		test_idp2.Id(): test_idp2,
		test_idp3.Id(): test_idp3,
	}
	for _, info := range buff {
		idp := m[info.Issuer]
		if idp == nil {
			t.Fatal("no ID provider " + info.Issuer)
		} else if info.Issuer_name != idp.Names()[""] {
			t.Error(info.Issuer_name)
			t.Fatal(idp.Names()[""])
		} else if info.Issuer_nameJa != idp.Names()["ja"] {
			t.Error(info.Issuer_nameJa)
			t.Fatal(idp.Names()["ja"])
		}
	}
}
