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
	"testing"
)

func TestSelectRequest(t *testing.T) {
	r, err := http.NewRequest("GET", "https://selector.example.org/select?ticket="+url.QueryEscape(test_ticId)+"&issuer="+url.QueryEscape(test_idp1.Id())+"&locale="+url.QueryEscape(test_lang), nil)
	if err != nil {
		t.Fatal(err)
	}

	req, err := parseSelectRequest(r)
	if err != nil {
		t.Fatal(err)
	} else if req.ticket() != test_ticId {
		t.Error(req.ticket())
		t.Fatal(test_ticId)
	} else if req.idProvider() != test_idp1.Id() {
		t.Error(req.idProvider())
		t.Fatal(test_idp1.Id())
	} else if req.language() != test_lang {
		t.Error(req.language())
		t.Fatal(test_lang)
	}
}
