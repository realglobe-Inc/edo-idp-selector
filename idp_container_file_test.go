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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const filePerm = 0644

func TestFileIdpContainer(t *testing.T) {
	path, err := ioutil.TempDir("", "edo-idp-selector")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	idpCont := newFileIdpContainer(path, 0, 0)
	for _, idp := range []*idProvider{testIdp, testIdp2} {
		idpPath := filepath.Join(path, keyToEscapedJsonPath(idp.Id))
		buff, err := json.Marshal(idp)
		if err != nil {
			t.Fatal(err)
		}
		if err := ioutil.WriteFile(idpPath, buff, filePerm); err != nil {
			t.Fatal(err)
		}
	}
	testIdpContainer(t, idpCont)
}
