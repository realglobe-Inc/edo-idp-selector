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
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestIdProviderRequest(t *testing.T) {
	r, err := http.NewRequest("GET", "https://selector.example.org/api/info/issuer?issuer="+url.QueryEscape("a+b*"), nil)
	if err != nil {
		t.Fatal(err)
	}

	req, err := parseIdProviderRequest(r)
	if err != nil {
		t.Fatal(err)
	} else if filt := map[string]string{"issuer": "a+b*"}; !reflect.DeepEqual(req.filter(), filt) {
		t.Error(req.filter())
		t.Fatal(filt)
	}
}
