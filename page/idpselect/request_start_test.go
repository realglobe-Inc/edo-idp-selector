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
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestStartRequest(t *testing.T) {
	query := "prompt=" + url.QueryEscape("select_account login") +
		"&display=page&ui_locales=" + url.QueryEscape("ja-JP en-US")
	r, err := http.NewRequest("GET", "https://selector.example.org/?"+query, nil)
	if err != nil {
		t.Fatal(err)
	}

	req, err := parseStartRequest(r)
	if err != nil {
		t.Fatal(err)
	} else if req.query() != query {
		t.Error(req.query())
		t.Fatal(query)
	} else if !req.selectForced() {
		t.Fatal("account selection is not forced")
	} else if req.display() != "page" {
		t.Error(req.display())
		t.Fatal("page")
	} else if langs := []string{"ja-JP", "en-US"}; !reflect.DeepEqual(req.languages(), langs) {
		t.Error(req.languages())
		t.Fatal(langs)
	}
}
