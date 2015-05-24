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

package idpselect

import (
	idpdb "github.com/realglobe-Inc/edo-idp-selector/database/idp"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestStartPage(t *testing.T) {
	page := newTestPage([]idpdb.Element{
		test_idp1,
	}, nil)

	r, err := http.NewRequest("GET", "https://selector.example.org/?"+test_query, nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	page.HandleStart(w, r)

	if w.Code != http.StatusFound {
		t.Error(w.Code)
		t.Fatal(http.StatusFound)
	} else if uri, err := url.Parse(w.HeaderMap.Get("Location")); err != nil {
		t.Fatal(err)
	} else if uri.Path != page.pathSelUi {
		t.Error(uri.Path)
		t.Fatal(page.pathSelUi)
	}
}
