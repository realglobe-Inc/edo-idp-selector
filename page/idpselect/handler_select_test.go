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
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestSelectPage(t *testing.T) {
	page := newTestPage([]idpdb.Element{
		test_idp1,
	}, nil)

	sess := session.New(test_sessId, time.Now().Add(page.sessExpIn))
	sess.SetQuery(test_query)
	sess.SetTicket(session.NewTicket(test_ticId, time.Now().Add(page.ticExpIn)))
	page.sessDb.Save(sess, time.Now().Add(page.sessDbExpIn))

	r, err := http.NewRequest("GET", "https://selector.example.org/select?ticket="+url.QueryEscape(sess.Ticket().Id())+"&issuer="+url.QueryEscape(test_idp1.Id()), nil)
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
