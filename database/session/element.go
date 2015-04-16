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
	req string
	// 現在発行されているチケット。
	tic string
	// 過去に選択された IdP の ID。
	pastIdps list.List
	// 最後に選択された表示言語。
	lang string
}

func New(id string, exp time.Time) *Element {
	return &Element{
		id:  id,
		exp: exp,
	}
}

// 履歴を引き継いだセッションを作成する。
func (this *Element) New(id string, exp time.Time) *Element {
	elem := &Element{
		id:   id,
		exp:  exp,
		lang: this.lang,
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
func (this *Element) ExpiresIn() time.Time {
	return this.exp
}

// 最後に選択された IdP の ID を返す。
func (this *Element) IdProvider() string {
	return this.idp
}

// IdP が選択されたことを反映させる。
func (this *Element) SelectIdProvider(idp string) {
	if this.idp == idp {
		return
	} else {
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
		if elem.Value == idp {
			this.pastIdps.Remove(elem)
			return
		}
	}
}

// 現在のリクエスト内容を返す。
func (this *Element) Request() string {
	return this.req
}

// リクエスト内容を保存する。
// URL のクエリ部分を想定。
func (this *Element) SetRequest(req string) {
	this.req = req
}

// 現在発行されているチケットを返す。
func (this *Element) Ticket() string {
	return this.tic
}

// チケットを保存する。
func (this *Element) SetTicket(tic string) {
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
	this.idp = ""
	this.req = ""
	this.tic = ""
}
