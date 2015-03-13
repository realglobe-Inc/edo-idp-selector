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
	"reflect"
	"testing"
)

var testIdp = &idProvider{
	Id:      "https://example.com",
	Name:    "sample idp",
	AuthUri: "https://example.com/login",
}
var testIdp2 = &idProvider{
	Id:      "idp-no-id",
	Name:    "認証装置2",
	AuthUri: "https://a.b.c.example.com/",
}

func testIdpContainer(t *testing.T, idpCont idpContainer) {
	defer idpCont.close()

	if idp, err := idpCont.get(testIdp.Id); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(idp, testIdp) {
		t.Error(idp)
	}

	if idps, err := idpCont.list(nil); err != nil {
		t.Fatal(err)
	} else if len(idps) != 2 {
		t.Error(idps)
	}
}
