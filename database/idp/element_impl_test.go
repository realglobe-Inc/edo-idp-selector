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
	"encoding/json"
	webdb "github.com/realglobe-Inc/edo-idp-selector/database/web"
	"testing"
)

func TestElementImpl(t *testing.T) {
	testElement(t, newElement(test_id, test_names, test_authUri, test_tokUri, test_acntUri, test_coopFrUri, test_keys))
}

func TestElementImplKeyDownload(t *testing.T) {
	keyUri := "https://example.org/keys"
	data, _ := json.Marshal([]interface{}{test_key.ToMap()})
	webDb := webdb.NewMemoryDb([]webdb.Element{
		webdb.New(keyUri, data),
	})

	elem := newElement(test_id, test_names, test_authUri, test_tokUri, test_acntUri, test_coopFrUri, nil)
	elem.keyUri = keyUri
	elem.setWebDbIfNeeded(webDb)
	testElement(t, elem)
}