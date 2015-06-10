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
	"github.com/realglobe-Inc/edo-idp-selector/request"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"
)

// 正常に選択 UI にリダイレクトされることの検査。
// セッションが無ければセッションが発行されることの検査。
func TestStartPage(t *testing.T) {
	page := newTestPage([]idpdb.Element{test_idp1}, nil)

	r, err := http.NewRequest("GET", "https://selector.example.org/?"+test_query, nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	page.HandleStart(w, r)

	if w.Code != http.StatusFound {
		t.Error(w.Code)
		t.Fatal(http.StatusFound)
	} else if uri, err := url.Parse(w.HeaderMap.Get("Location")); err != nil {
		t.Fatal(err)
	} else if uri.Path != page.pathSelUi {
		t.Error(uri.Path)
		t.Fatal(page.pathSelUi)
	} else if ok, err := regexp.MatchString(page.sessLabel+"=[0-9a-zA-Z_\\-]", w.HeaderMap.Get("Set-Cookie")); err != nil {
		t.Fatal(err)
	} else if !ok {
		t.Error("no new session")
		t.Fatal(w.HeaderMap.Get("Set-Cookie"))
	}
}

// セッションが有効ならセッションが発行されないことの検査。
func TestStartPageNoSessionPublication(t *testing.T) {
	page := newTestPage([]idpdb.Element{test_idp1}, nil)
	now := time.Now()
	page.sessDb.Save(session.New(test_sessId, now.Add(page.sessExpIn)), now.Add(page.sessDbExpIn))

	r, err := http.NewRequest("GET", "https://selector.example.org/?"+test_query, nil)
	if err != nil {
		t.Fatal(err)
	}
	r.AddCookie(&http.Cookie{
		Name:  page.sessLabel,
		Value: test_sessId,
	})

	w := httptest.NewRecorder()
	page.HandleStart(w, r)

	if ok, err := regexp.MatchString(page.sessLabel+"=[0-9a-zA-Z_\\-]", w.HeaderMap.Get("Set-Cookie")); err != nil {
		t.Fatal(err)
	} else if ok {
		t.Error("new session")
		t.Fatal(w.HeaderMap.Get("Set-Cookie"))
	}
}

// セッションの期限が切れそうならセッションが更新されることの検査。
func TestStartPageSessionRefresh(t *testing.T) {
	page := newTestPage([]idpdb.Element{test_idp1}, nil)
	now := time.Now()
	sess := session.New(test_sessId, now.Add(page.sessRefDelay-time.Nanosecond))
	sess.SelectIdProvider(test_idp1.Id())
	page.sessDb.Save(sess, now.Add(page.sessDbExpIn))

	r, err := http.NewRequest("GET", "https://selector.example.org/?"+test_query, nil)
	if err != nil {
		t.Fatal(err)
	}
	r.AddCookie(&http.Cookie{
		Name:  page.sessLabel,
		Value: test_sessId,
	})

	w := httptest.NewRecorder()
	page.HandleStart(w, r)

	if ok, err := regexp.MatchString(page.sessLabel+"=[0-9a-zA-Z_\\-]", w.HeaderMap.Get("Set-Cookie")); err != nil {
		t.Fatal(err)
	} else if !ok {
		t.Error("no new session")
		t.Fatal(w.HeaderMap.Get("Set-Cookie"))
	} else if uri, err := url.Parse(w.HeaderMap.Get("Location")); err != nil {
		t.Fatal(err)
	} else if uri.Path != page.pathSelUi {
		t.Error(uri.Path)
		t.Fatal(page.pathSelUi)
	} else if query := uri.Query(); len(query) == 0 {
		t.Fatal("no query")
	} else if isss, isss2 := `["`+test_idp1.Id()+`"]`, query.Get("issuers"); isss2 != isss {
		t.Error(isss2)
		t.Fatal(isss)
	}
}

// ID プロバイダと紐付くセッションなら ID プロバイダにリダイレクトされることの検査。
func TestStartPageRedirectIdProvider(t *testing.T) {
	page := newTestPage([]idpdb.Element{test_idp1}, nil)
	now := time.Now()
	sess := session.New(test_sessId, now.Add(page.sessExpIn))
	sess.SelectIdProvider(test_idp1.Id())
	page.sessDb.Save(sess, now.Add(page.sessDbExpIn))

	r, err := http.NewRequest("GET", "https://selector.example.org/?"+test_query, nil)
	if err != nil {
		t.Fatal(err)
	}
	r.AddCookie(&http.Cookie{
		Name:  page.sessLabel,
		Value: test_sessId,
	})

	w := httptest.NewRecorder()
	page.HandleStart(w, r)

	if w.Code != http.StatusFound {
		t.Error(w.Code)
		t.Fatal(http.StatusFound)
	} else if w.HeaderMap.Get("Location") != test_idp1.AuthUri()+"?"+test_query {
		t.Error(w.HeaderMap.Get("Location"))
		t.Fatal(test_idp1.AuthUri() + "?" + test_query)
	}
}

// ID プロバイダと紐付くセッションでも select_account パラメータ付きなら
// 選択 UI にリダイレクトされることの検査。
func TestStartPageRedirectUi(t *testing.T) {
	page := newTestPage([]idpdb.Element{test_idp1}, nil)
	now := time.Now()
	sess := session.New(test_sessId, now.Add(page.sessExpIn))
	sess.SelectIdProvider(test_idp1.Id())
	page.sessDb.Save(sess, now.Add(page.sessDbExpIn))

	var query string
	{
		q, err := url.ParseQuery(test_query)
		if err != nil {
			t.Fatal(err)
		}
		prmpt := request.FormValueSet(q.Get("prompt"))
		prmpt["select_account"] = true
		q.Set("prompt", request.ValueSetForm(prmpt))
		query = q.Encode()
	}

	r, err := http.NewRequest("GET", "https://selector.example.org/?"+query, nil)
	if err != nil {
		t.Fatal(err)
	}
	r.AddCookie(&http.Cookie{
		Name:  page.sessLabel,
		Value: test_sessId,
	})

	w := httptest.NewRecorder()
	page.HandleStart(w, r)

	if w.Code != http.StatusFound {
		t.Error(w.Code)
		t.Fatal(http.StatusFound)
	} else if uri, err := url.Parse(w.HeaderMap.Get("Location")); err != nil {
		t.Fatal(err)
	} else if uri.Path != page.pathSelUi {
		t.Error(uri.Path)
		t.Fatal(page.pathSelUi)
	}
}
