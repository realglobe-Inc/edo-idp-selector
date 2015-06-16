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
	"container/list"
	"encoding/json"
	"github.com/realglobe-Inc/edo-idp-selector/ticket"
	rist "github.com/realglobe-Inc/edo-lib/list"
	"github.com/realglobe-Inc/go-lib/erro"
	"time"
)

// セッション。
type Element struct {
	id string
	// 有効期限。
	exp time.Time
	// 最後に選択された IdP の ID。
	idp string
	// 現在のリクエスト内容。
	query string
	// 現在発行されているチケット。
	tic *ticket.Ticket
	// 過去に選択された IdP の ID。
	pastIdps *list.List
	// 最後に選択された表示言語。
	lang string

	// 以下、作業用。

	// 読み込まれたセッションかどうか。
	saved bool
}

// 防御的コピー用。
func (this *Element) copy() *Element {
	elem := New(this.id, this.exp)
	elem.idp = this.idp
	elem.query = this.query
	elem.tic = this.tic
	for e := this.pastIdps.Back(); e != nil; e = e.Prev() {
		elem.pastIdps.PushFront(e.Value)
	}
	elem.lang = this.lang
	return elem
}

func New(id string, exp time.Time) *Element {
	return &Element{
		id:       id,
		exp:      exp,
		pastIdps: list.New(),
	}
}

// 履歴を引き継いだセッションを作成する。
func (this *Element) New(id string, exp time.Time) *Element {
	elem := &Element{
		id:       id,
		exp:      exp,
		lang:     this.lang,
		pastIdps: list.New(),
	}
	for e := this.pastIdps.Back(); e != nil; e = e.Prev() {
		elem.pastIdps.PushFront(e.Value)
	}
	if this.idp != "" {
		elem.addPastIdProvider(this.idp, MaxHistory)
	}
	return elem
}

// 過去に選択された IdP をいくつまで記憶するか。
// 最後に選択された IdP も含む。
var MaxHistory = 5

// ID を返す。
func (this *Element) Id() string {
	return this.id
}

// 有効期限を返す。
func (this *Element) Expires() time.Time {
	return this.exp
}

// 最後に選択された IdP の ID を返す。
func (this *Element) IdProvider() string {
	return this.idp
}

// IdP が選択されたことを反映させる。
func (this *Element) SelectIdProvider(idp string) {
	if this.idp != idp {
		this.removePastIdProvider(idp)
		if this.idp != "" {
			this.addPastIdProvider(this.idp, MaxHistory-1)
		}
	}
	this.idp = idp
}

func (this *Element) addPastIdProvider(idp string, max int) {
	for this.pastIdps.Len() >= max {
		this.pastIdps.Remove(this.pastIdps.Back())
	}
	this.pastIdps.PushFront(idp)
	return
}

func (this *Element) removePastIdProvider(idp string) {
	for elem := this.pastIdps.Front(); elem != nil; elem = elem.Next() {
		if elem.Value.(string) == idp {
			this.pastIdps.Remove(elem)
			return
		}
	}
}

// 現在のリクエスト内容を返す。
func (this *Element) Query() string {
	return this.query
}

// リクエスト内容を保存する。
func (this *Element) SetQuery(query string) {
	this.query = query
}

// 現在発行されているチケットを返す。
func (this *Element) Ticket() *ticket.Ticket {
	return this.tic
}

// チケットを保存する。
func (this *Element) SetTicket(tic *ticket.Ticket) {
	this.tic = tic
}

// 過去に選択された IdP の ID を返す。
func (this *Element) SelectedIdProviders() []string {
	a := []string{}
	if this.idp != "" {
		a = append(a, this.idp)
	}
	for elem := this.pastIdps.Front(); elem != nil; elem = elem.Next() {
		a = append(a, elem.Value.(string))
	}
	return a
}

// 最後に選択された表示言語を返す。
func (this *Element) Language() string {
	return this.lang
}

// 表示言語を保存する。
func (this *Element) SetLanguage(lang string) {
	this.lang = lang
}

// 一時データを消す。
func (this *Element) Clear() {
	this.query = ""
	this.tic = nil
}

// 読み込まれたセッションかどうか。
func (this *Element) Saved() bool {
	return this.saved
}

func (this *Element) setSaved() {
	this.saved = true
}

//  {
//      "id": <ID>,
//      "expires": <有効期限>,
//      "issuer": <ID プロバイダ>,
//      "query": <リクエスト内容>,
//      "ticket": <チケット>,
//      "past_issuers": [
//          <選択したことのある ID プロバイダ>,
//          ...
//      ],
//      "locale": <表示言語>
//  }
func (this *Element) MarshalJSON() (data []byte, err error) {
	return json.Marshal(map[string]interface{}{
		"id":           this.id,
		"expires":      this.exp,
		"issuer":       this.idp,
		"query":        this.query,
		"ticket":       this.tic,
		"past_issuers": (*rist.List)(this.pastIdps),
		"locale":       this.lang,
	})
}

func (this *Element) UnmarshalJSON(data []byte) error {
	var buff struct {
		Id       string         `json:"id"`
		Exp      time.Time      `json:"expires"`
		Idp      string         `json:"issuer"`
		Query    string         `json:"query"`
		Tic      *ticket.Ticket `json:"ticket"`
		PastIdps *rist.List     `json:"past_issuers"`
		Lang     string         `json:"locale"`
	}
	if err := json.Unmarshal(data, &buff); err != nil {
		return erro.Wrap(err)
	}

	// 中身を文字列に限定。
	pastIdps := (*list.List)(buff.PastIdps)
	filted := list.New()
	for e := pastIdps.Back(); e != nil; e = e.Prev() {
		idp, ok := e.Value.(string)
		if ok {
			filted.PushFront(idp)
		}
	}

	this.id = buff.Id
	this.exp = buff.Exp
	this.idp = buff.Idp
	this.query = buff.Query
	this.tic = buff.Tic
	this.pastIdps = filted
	this.lang = buff.Lang
	return nil
}
