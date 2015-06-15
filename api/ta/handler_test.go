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

package ta

import (
	"encoding/json"
	tadb "github.com/realglobe-Inc/edo-idp-selector/database/ta"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHandler(t *testing.T) {
	ta := tadb.New("https://ta.example.org",
		map[string]string{"": "The TA", "ja": "かの TA"},
		map[string]bool{"https://ta.example.org/callback": true},
		nil,
		true,
		"https://ta.example.org")
	db := tadb.NewMemoryDb([]tadb.Element{ta})
	hndl := New(nil, "/api/info/ta", db, true)

	r, err := http.NewRequest("GET", "https://example.org/api/info/ta/"+url.QueryEscape(url.QueryEscape(ta.Id())), nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	hndl.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Error(w.Code)
		t.Fatal(http.StatusOK)
	} else if w.HeaderMap.Get("Content-Type") != "application/json" {
		t.Error(w.HeaderMap.Get("Content-Type"))
		t.Fatal("application/json")
	} else if w.Body == nil {
		t.Fatal("no body")
	}

	data, _ := ioutil.ReadAll(w.Body)
	var buff struct {
		Client_name    string
		Client_name_ja string `json:"client_name#ja"`
	}
	if err := json.Unmarshal(data, &buff); err != nil {
		t.Fatal(err)
	} else if buff.Client_name != "The TA" {
		t.Error(buff.Client_name)
		t.Fatal("The TA")
	} else if buff.Client_name_ja != "かの TA" {
		t.Error(buff.Client_name_ja)
		t.Fatal("かの TA")
	}
}
