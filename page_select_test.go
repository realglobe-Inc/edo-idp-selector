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
	idpdb "github.com/realglobe-Inc/edo-idp-selector/database/idp"
	"github.com/realglobe-Inc/edo-idp-selector/database/session"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestSelectPage(t *testing.T) {
	sys := newTestSystem(nil, []idpdb.Element{
		test_idp1,
	}, nil)

	sess := session.New(test_sessId, time.Now().Add(sys.sessExpIn))
	sess.SetQuery(test_query)
	sess.SetTicket(session.NewTicket(test_ticId, time.Now().Add(sys.ticExpIn)))
	sys.sessDb.Save(sess, time.Now().Add(sys.sessDbExpIn))

	r, err := http.NewRequest("GET", "https://selector.example.org/select?ticket="+url.QueryEscape(sess.Ticket().Id())+"&issuer="+url.QueryEscape(test_idp1.Id()), nil)
	if err != nil {
		t.Fatal(err)
	}
	r.AddCookie(&http.Cookie{
		Name:  sys.sessLabel,
		Value: sess.Id(),
	})

	w := httptest.NewRecorder()
	if err := sys.selectPage(w, r); err != nil {
		t.Fatal(err)
	} else if w.Code != http.StatusFound {
		t.Error(w.Code)
		t.Fatal(http.StatusFound)
	} else if w.HeaderMap.Get("Location") != test_idp1.AuthenticationUri()+"?"+test_query {
		t.Error(w.HeaderMap.Get("Location"))
		t.Fatal(test_idp1.AuthenticationUri() + "?" + test_query)
	}
}
