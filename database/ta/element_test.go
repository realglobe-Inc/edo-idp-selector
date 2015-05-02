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
	"github.com/realglobe-Inc/edo-lib/jwk"
	"reflect"
	"testing"
)

const (
	test_id      = "https://ta.example.org"
	test_rediUri = "https://ta.example.org/return"
	test_pw      = true
	test_sect    = "https://ta.example.org"
)

var (
	test_names    = map[string]string{"": "The TA", "ja": "かの TA"}
	test_rediUris = map[string]bool{test_rediUri: true}
	test_key, _   = jwk.FromMap(map[string]interface{}{
		"kty": "EC",
		"crv": "P-256",
		"x":   "lpHYO1qpjU95B2sThPR2-1jv44axgaEDkQtcKNE-oZs",
		"y":   "soy5O11SFFFeYdhQVodXlYPIpeo0pCS69IxiVPPf0Tk",
		"d":   "3BhkCluOkm8d8gvaPD5FDG2zeEw2JKf3D5LwN-mYmsw",
	})
	test_keys = []jwk.Key{test_key}
)

func testElement(t *testing.T, elem Element) {
	if elem.Id() != test_id {
		t.Error(elem.Id())
		t.Fatal(test_id)
	} else if !reflect.DeepEqual(elem.Names(), test_names) {
		t.Error(elem.Names())
		t.Fatal(test_names)
	} else if !reflect.DeepEqual(elem.RedirectUris(), test_rediUris) {
		t.Error(elem.RedirectUris())
		t.Fatal(test_rediUris)
	} else if !reflect.DeepEqual(elem.Keys(), test_keys) {
		t.Error(elem.Keys())
		t.Fatal(test_keys)
	} else if elem.Pairwise() != test_pw {
		t.Error(elem.Pairwise())
		t.Fatal(test_pw)
	} else if elem.Sector() != test_sect {
		t.Error(elem.Sector())
		t.Fatal(test_sect)
	}
}
