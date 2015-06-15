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

package idpselect

import (
	idpdb "github.com/realglobe-Inc/edo-idp-selector/database/idp"
	"github.com/realglobe-Inc/edo-idp-selector/database/session"
	tadb "github.com/realglobe-Inc/edo-idp-selector/database/ta"
	"github.com/realglobe-Inc/edo-idp-selector/ticket"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"
)

// 正常に ID プロバイダにリダイレクトされることの検査。
func TestSelectPage(t *testing.T) {
	page := newTestPage([]idpdb.Element{
		test_idp1,
	}, nil)

	now := time.Now()
	sess := session.New(test_sessId, now.Add(page.sessExpIn))
	sess.SetQuery(test_query)
	sess.SetTicket(ticket.New(test_ticId, now.Add(page.ticExpIn)))
	page.sessDb.Save(sess, now.Add(page.sessDbExpIn))

	r, err := http.NewRequest("GET", "https://selector.example.org/select"+
		"?ticket="+url.QueryEscape(sess.Ticket().Id())+
		"&issuer="+url.QueryEscape(test_idp1.Id()), nil)
	if err != nil {
		t.Fatal(err)
	}
	r.AddCookie(&http.Cookie{
		Name:  page.sessLabel,
		Value: sess.Id(),
	})

	w := httptest.NewRecorder()
	page.HandleSelect(w, r)

	if w.Code != http.StatusFound {
		t.Error(w.Code)
		t.Fatal(http.StatusFound)
	} else if w.HeaderMap.Get("Location") != test_idp1.AuthUri()+"?"+test_query {
		t.Error(w.HeaderMap.Get("Location"))
		t.Fatal(test_idp1.AuthUri() + "?" + test_query)
	}
}

// セッションが無ければセッションが発行されることの検査。
func TestSelectPageSessionPublication(t *testing.T) {
	page := newTestPage([]idpdb.Element{
		test_idp1,
	}, nil)

	r, err := http.NewRequest("GET", "https://selector.example.org/select"+
		"?ticket="+url.QueryEscape(test_ticId)+
		"&issuer="+url.QueryEscape(test_idp1.Id()), nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	page.HandleSelect(w, r)

	if w.Code == http.StatusFound {
		t.Fatal("not error")
	} else if ok, err := regexp.MatchString(page.sessLabel+"=[0-9a-zA-Z_\\-]", w.HeaderMap.Get("Set-Cookie")); err != nil {
		t.Fatal(err)
	} else if !ok {
		t.Error("no new session")
		t.Fatal(w.HeaderMap.Get("Set-Cookie"))
	}
}

// セッションが有効ならセッションが発行されないことの検査。
func TestSelectPageNoSessionPublication(t *testing.T) {
	page := newTestPage([]idpdb.Element{
		test_idp1,
	}, nil)

	now := time.Now()
	sess := session.New(test_sessId, now.Add(page.sessExpIn))
	sess.SetQuery(test_query)
	sess.SetTicket(ticket.New(test_ticId, now.Add(page.ticExpIn)))
	page.sessDb.Save(sess, now.Add(page.sessDbExpIn))

	r, err := http.NewRequest("GET", "https://selector.example.org/select"+
		"?ticket="+url.QueryEscape(sess.Ticket().Id())+
		"&issuer="+url.QueryEscape(test_idp1.Id()), nil)
	if err != nil {
		t.Fatal(err)
	}
	r.AddCookie(&http.Cookie{
		Name:  page.sessLabel,
		Value: sess.Id(),
	})

	w := httptest.NewRecorder()
	page.HandleSelect(w, r)

	if ok, err := regexp.MatchString(page.sessLabel+"=[0-9a-zA-Z_\\-]", w.HeaderMap.Get("Set-Cookie")); err != nil {
		t.Fatal(err)
	} else if ok {
		t.Error("new session")
		t.Fatal(w.HeaderMap.Get("Set-Cookie"))
	}
}

// チケットが無効ならエラーリダイレクトされることの検査。
func TestSelectPageErrorRedirect(t *testing.T) {
	page := newTestPage([]idpdb.Element{
		test_idp1,
	}, []tadb.Element{
		test_ta,
	})

	now := time.Now()
	sess := session.New(test_sessId, now.Add(page.sessExpIn))
	sess.SetQuery(test_query)
	sess.SetTicket(ticket.New(test_ticId, now.Add(page.ticExpIn)))
	page.sessDb.Save(sess, now.Add(page.sessDbExpIn))

	r, err := http.NewRequest("GET", "https://selector.example.org/select"+
		"?ticket="+url.QueryEscape(sess.Ticket().Id()+"a")+
		"&issuer="+url.QueryEscape(test_idp1.Id()), nil)
	if err != nil {
		t.Fatal(err)
	}
	r.AddCookie(&http.Cookie{
		Name:  page.sessLabel,
		Value: sess.Id(),
	})

	w := httptest.NewRecorder()
	page.HandleSelect(w, r)

	if w.Code != http.StatusFound {
		t.Error(w.Code)
		t.Fatal(http.StatusFound)
	} else if uri, err := url.Parse(w.HeaderMap.Get("Location")); err != nil {
		t.Fatal(err)
	} else if rediUri := uri.Scheme + "://" + uri.Host + uri.Path; !test_ta.RedirectUris()[rediUri] {
		t.Error("not redirect uri")
		t.Error(rediUri)
		t.Fatal(test_ta.RedirectUris())
	} else if uri.Query().Get("error") == "" {
		t.Fatal("no error")
	}
}
