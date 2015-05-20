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
	"github.com/realglobe-Inc/go-lib/erro"
	"time"
)

// web API を叩くための許可証。
type Ticket struct {
	id string
	// 有効期限。
	exp time.Time
}

func NewTicket(id string, exp time.Time) *Ticket {
	return &Ticket{
		id:  id,
		exp: exp,
	}
}

// ID を返す。
func (this *Ticket) Id() string {
	return this.id
}

// 有効期限を返す。
func (this *Ticket) Expires() time.Time {
	return this.exp
}

//  {
//      "id": <ID>,
//      "expires": <有効期限>
//  }
func (this *Ticket) MarshalJSON() (data []byte, err error) {
	m := map[string]interface{}{
		"id":      this.id,
		"expires": this.exp,
	}
	return json.Marshal(m)
}

func (this *Ticket) UnmarshalJSON(data []byte) error {
	var buff struct {
		Id  string    `json:"id"`
		Exp time.Time `json:"expires"`
	}
	if err := json.Unmarshal(data, &buff); err != nil {
		return erro.Wrap(err)
	}
	this.id = buff.Id
	this.exp = buff.Exp
	return nil
}
