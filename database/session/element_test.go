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
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestElementPastIdProvider(t *testing.T) {
	a := New("test-session-id", time.Date(2015, time.April, 4, 18, 41, 20, 123456789, time.UTC))
	if idps := a.SelectedIdProviders(); len(idps) != 0 {
		t.Error(idps)
		return
	}

	a.SelectIdProvider("test-id-provider")
	if idps := a.SelectedIdProviders(); len(idps) != 1 {
		t.Error(idps)
		return
	}

	a.SelectIdProvider("test-id-provider2")
	if idps := a.SelectedIdProviders(); len(idps) != 2 {
		t.Error(idps)
		return
	}

	a.SelectIdProvider("test-id-provider")
	if idps := a.SelectedIdProviders(); len(idps) != 2 {
		t.Error(idps)
		return
	}

	a.SelectIdProvider("test-id-provider3")
	if idps := a.SelectedIdProviders(); len(idps) != 3 {
		t.Error(idps)
		return
	}

	if idps := a.SelectedIdProviders(); !reflect.DeepEqual(idps, []string{"test-id-provider3", "test-id-provider", "test-id-provider2"}) {
		t.Error(idps)
		return
	}

	for i := 0; i < 2*MaxHistory; i++ {
		a.SelectIdProvider("test-id-provider" + strconv.Itoa(i))
		if idps := a.SelectedIdProviders(); len(idps) > MaxHistory+1 {
			t.Error(i)
			t.Error(idps)
			return
		}
	}
	if idps := a.SelectedIdProviders(); len(idps) != MaxHistory {
		t.Error(idps)
		return
	}
}

func TestElementNew(t *testing.T) {
	date := time.Date(2015, time.April, 4, 18, 41, 20, 123456789, time.UTC)
	a := New("test-session-id", date)
	a.SetRequest("param=val&param2=val2")
	a.SetTicket("test-ticket")
	a.SetLanguage("test-language")
	for i := 0; i < 2*MaxHistory; i++ {
		a.SelectIdProvider("test-id-provider" + strconv.Itoa(i))
		b := a.New("test-session-id2", date.Add(time.Second))

		if b.Id() == a.Id() {
			t.Error(i)
			t.Error(b.Id())
			return
		} else if b.ExpiresIn().Equal(a.ExpiresIn()) {
			t.Error(i)
			t.Error(b.ExpiresIn())
			return
		} else if b.IdProvider() != "" {
			t.Error(i)
			t.Error(b.IdProvider())
			return
		} else if b.Request() != "" {
			t.Error(i)
			t.Error(b.Request())
			return
		} else if b.Ticket() != "" {
			t.Error(i)
			t.Error(b.Ticket())
			return
		} else if idps, idps2 := a.SelectedIdProviders(), b.SelectedIdProviders(); !reflect.DeepEqual(idps, idps2) {
			t.Error(i)
			t.Error(idps2)
			t.Error(idps)
			return
		} else if b.Language() != a.Language() {
			t.Error(i)
			t.Error(b.Language())
			t.Error(a.Language())
			return
		}
	}
}
