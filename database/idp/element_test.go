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
	"reflect"
	"testing"

	"github.com/realglobe-Inc/edo-lib/jwk"
)

const (
	test_id        = "https://idp.example.org"
	test_name      = "The ID Provider"
	test_nameJa    = "かの ID プロバイダ"
	test_authUri   = "https://idp.example.org/auth"
	test_tokUri    = "https://idp.example.org/token"
	test_acntUri   = "https://idp.example.org/info/account"
	test_coopFrUri = "https://idp.example.org/coop/from"
	test_coopToUri = "https://idp.example.org/coop/to"
)

var (
	test_names  = map[string]string{"": test_name, "ja": test_nameJa}
	test_key, _ = jwk.FromMap(map[string]interface{}{
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
	} else if elem.AuthUri() != test_authUri {
		t.Error(elem.AuthUri())
		t.Fatal(test_authUri)
	} else if elem.TokenUri() != test_tokUri {
		t.Error(elem.TokenUri())
		t.Fatal(test_tokUri)
	} else if elem.AccountUri() != test_acntUri {
		t.Error(elem.AccountUri())
		t.Fatal(test_acntUri)
	} else if elem.CoopFromUri() != test_coopFrUri {
		t.Error(elem.CoopFromUri())
		t.Fatal(test_coopFrUri)
	} else if elem.CoopToUri() != test_coopToUri {
		t.Error(elem.CoopToUri())
		t.Fatal(test_coopToUri)
	} else if !reflect.DeepEqual(elem.Keys(), test_keys) {
		t.Error(elem.Keys())
		t.Fatal(test_keys)
	}
}
