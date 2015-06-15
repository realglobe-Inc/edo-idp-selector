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

package session

import (
	"encoding/json"
	"github.com/realglobe-Inc/edo-idp-selector/ticket"
	"reflect"
	"strconv"
	"testing"
	"time"
)

const (
	test_id    = "pbqy9Fx6OKqyGFe6FYS8SsqzZNuWxL"
	test_idp   = "https://idp.example.org"
	test_query = "response_type=code&scope=openid&client_id=https%3A%2F%2Fta.example.org&redirect_uri=https%3A%2F%2Fta.example.org%2Fcallback"
	test_lang  = "ja"
	test_ticId = "2IHV7qG7SJ"
)

var (
	test_tic = ticket.New(test_ticId, time.Now().Add(24*time.Hour))
)

func TestElement(t *testing.T) {
	exp := time.Now().Add(24 * time.Hour)
	elem := New(test_id, exp)

	if elem.Id() != test_id {
		t.Error(elem.Id())
		t.Fatal(test_id)
	} else if !elem.Expires().Equal(exp) {
		t.Error(elem.Expires())
		t.Fatal(exp)
	} else if elem.IdProvider() != "" {
		t.Fatal(elem.IdProvider())
	} else if elem.Query() != "" {
		t.Fatal(elem.Query())
	} else if elem.Ticket() != nil {
		t.Fatal(elem.Ticket())
	} else if len(elem.SelectedIdProviders()) > 0 {
		t.Fatal(elem.SelectedIdProviders())
	} else if elem.Language() != "" {
		t.Fatal(elem.Language())
	}

	elem.SelectIdProvider(test_idp)
	elem.SetQuery(test_query)
	elem.SetTicket(test_tic)
	elem.SetLanguage(test_lang)

	if elem.IdProvider() != test_idp {
		t.Error(elem.IdProvider())
		t.Fatal(test_idp)
	} else if elem.Query() != test_query {
		t.Error(elem.Query())
		t.Fatal(test_query)
	} else if !reflect.DeepEqual(elem.Ticket(), test_tic) {
		t.Error(elem.Ticket())
		t.Fatal(test_tic)
	} else if !reflect.DeepEqual(elem.SelectedIdProviders(), []string{test_idp}) {
		t.Error(elem.SelectedIdProviders())
		t.Fatal([]string{test_idp})
	} else if elem.Language() != test_lang {
		t.Error(elem.Language())
		t.Fatal(test_lang)
	}

	elem.Clear()
	if elem.Query() != "" {
		t.Fatal(elem.Query())
	} else if elem.Ticket() != nil {
		t.Fatal(elem.Ticket())
	}
}

func TestElementPastIdProvider(t *testing.T) {
	exp := time.Now().Add(24 * time.Hour)
	elem := New(test_id, exp)
	if len(elem.SelectedIdProviders()) != 0 {
		t.Fatal(elem.SelectedIdProviders())
	}

	elem.SelectIdProvider(test_idp)
	if len(elem.SelectedIdProviders()) != 1 {
		t.Fatal(elem.SelectedIdProviders())
	}

	elem.SelectIdProvider(test_idp + "2")
	if len(elem.SelectedIdProviders()) != 2 {
		t.Fatal(elem.SelectedIdProviders())
	}

	// 同じなら増えない。
	elem.SelectIdProvider(test_idp)
	if len(elem.SelectedIdProviders()) != 2 {
		t.Fatal(elem.SelectedIdProviders())
	}

	elem.SelectIdProvider(test_idp + "3")
	if len(elem.SelectedIdProviders()) != 3 {
		t.Fatal(elem.SelectedIdProviders())
	}

	if !reflect.DeepEqual(elem.SelectedIdProviders(), []string{
		test_idp + "3",
		test_idp,
		test_idp + "2"}) {
		t.Fatal(elem.SelectedIdProviders())
	}

	for i := 0; i < 2*MaxHistory; i++ {
		elem.SelectIdProvider(test_idp + strconv.Itoa(i))
		if len(elem.SelectedIdProviders()) > MaxHistory {
			t.Error(i)
			t.Fatal(elem.SelectedIdProviders())
		}
	}
	if len(elem.SelectedIdProviders()) != MaxHistory {
		t.Fatal(elem.SelectedIdProviders())
	}
}

func testElementJson(t *testing.T, elem *Element) {
	data, err := json.Marshal(elem)
	if err != nil {
		t.Fatal(err)
	}

	var elem2 Element
	if err := json.Unmarshal(data, &elem2); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&elem2, elem) {
		t.Error(&elem2)
		t.Fatal(elem)
	}
}

func TestElementJson(t *testing.T) {
	exp := time.Now().Add(24 * time.Hour)
	elem := New(test_id, exp)
	testElementJson(t, elem)

	elem.SelectIdProvider(test_idp)
	testElementJson(t, elem)

	elem.SetQuery(test_query)
	testElementJson(t, elem)

	elem.SetTicket(test_tic)
	testElementJson(t, elem)

	elem.SetLanguage(test_lang)
	testElementJson(t, elem)
}
