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

package ticket

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

const (
	test_ticId = "2IHV7qG7SJ"
)

func TestTicket(t *testing.T) {
	exp := time.Now().Add(24 * time.Hour)
	tic := New(test_ticId, exp)

	if tic.Id() != test_ticId {
		t.Error(tic.Id())
		t.Fatal(test_ticId)
	} else if !tic.Expires().Equal(exp) {
		t.Error(tic.Expires())
		t.Fatal(exp)
	}
}

func TestTicketJson(t *testing.T) {
	exp := time.Now().Add(24 * time.Hour)
	tic := New(test_ticId, exp)

	data, err := json.Marshal(tic)
	if err != nil {
		t.Fatal(err)
	}

	var tic2 Ticket
	if err := json.Unmarshal(data, &tic2); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&tic2, tic) {
		t.Error(&tic2)
		t.Fatal(tic)
	}
}
