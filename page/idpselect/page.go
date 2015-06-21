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

// ID プロバイダ選択ページ。
package idpselect

import (
	idpdb "github.com/realglobe-Inc/edo-idp-selector/database/idp"
	"github.com/realglobe-Inc/edo-idp-selector/database/session"
	tadb "github.com/realglobe-Inc/edo-idp-selector/database/ta"
	idperr "github.com/realglobe-Inc/edo-idp-selector/error"
	"github.com/realglobe-Inc/edo-idp-selector/request"
	logutil "github.com/realglobe-Inc/edo-lib/log"
	"github.com/realglobe-Inc/edo-lib/rand"
	"github.com/realglobe-Inc/edo-lib/server"
	"github.com/realglobe-Inc/go-lib/erro"
	"html/template"
	"net/http"
	"net/url"
	"time"
)

type Page struct {
	stopper *server.Stopper

	pathSelUi string
	errTmpl   *template.Template

	sessLabel    string
	sessLen      int
	sessExpIn    time.Duration
	sessRefDelay time.Duration
	sessDbExpIn  time.Duration
	ticLen       int
	ticExpIn     time.Duration

	idpDb  idpdb.Db
	taDb   tadb.Db
	sessDb session.Db
	idGen  rand.Generator

	cookPath string
	cookSec  bool
	debug    bool
}

func New(
	stopper *server.Stopper,
	pathSelUi string,
	errTmpl *template.Template,
	sessLabel string,
	sessLen int,
	sessExpIn time.Duration,
	sessRefDelay time.Duration,
	sessDbExpIn time.Duration,
	ticLen int,
	ticExpIn time.Duration,
	idpDb idpdb.Db,
	taDb tadb.Db,
	sessDb session.Db,
	idGen rand.Generator,
	cookPath string,
	cookSec bool,
	debug bool,
) *Page {
	return &Page{
		stopper,
		pathSelUi,
		errTmpl,
		sessLabel,
		sessLen,
		sessExpIn,
		sessRefDelay,
		sessDbExpIn,
		ticLen,
		ticExpIn,
		idpDb,
		taDb,
		sessDb,
		idGen,
		cookPath,
		cookSec,
		debug,
	}
}

func (this *Page) newCookie(sess *session.Element) *http.Cookie {
	return &http.Cookie{
		Name:     this.sessLabel,
		Value:    sess.Id(),
		Path:     this.cookPath,
		Expires:  sess.Expires(),
		Secure:   this.cookSec,
		HttpOnly: true,
	}
}

// environment のメソッドは idperr.Error を返す。
type environment struct {
	*Page

	sender *request.Request
	sess   *session.Element
}

func (this *environment) respondErrorHtml(w http.ResponseWriter, r *http.Request, origErr error) {
	var uri *url.URL
	if this.sess.Query() != "" {
		var err error
		uri, err = getRedirectUri(this.sess.Query(), this.taDb)
		if err != nil {
			log.Err(this.sender, ": ", erro.Unwrap(err))
			log.Debug(this.sender, ": ", erro.Wrap(err))
		}
	}

	// 経過を破棄。
	this.sess.Clear()
	if err := this.sessDb.Save(this.sess, this.sess.Expires().Add(this.sessDbExpIn-this.sessExpIn)); err != nil {
		log.Err(this.sender, ": ", erro.Wrap(err))
	} else {
		log.Debug(this.sender, ": Saved session "+logutil.Mosaic(this.sess.Id()))
	}

	if !this.sess.Saved() {
		// 未通知セッションの通知。
		http.SetCookie(w, this.newCookie(this.sess))
		log.Debug(this.sender, ": Report session "+logutil.Mosaic(this.sess.Id()))
	}

	if uri != nil {
		idperr.RedirectError(w, r, origErr, uri, this.sender)
		return
	}

	idperr.RespondHtml(w, r, origErr, this.errTmpl, this.sender)
	return
}

// リクエストからリダイレクト URI を取得する。
func getRedirectUri(rawQuery string, taDb tadb.Db) (*url.URL, error) {
	vals, err := url.ParseQuery(rawQuery)
	if err != nil {
		return nil, erro.Wrap(err)
	} else if taId := vals.Get(tagClient_id); taId == "" {
		return nil, erro.New("no TA ID")
	} else if rawRediUri := vals.Get(tagRedirect_uri); rawRediUri == "" {
		return nil, erro.New("no redirect URI")
	} else if ta, err := taDb.Get(taId); err != nil {
		return nil, erro.Wrap(err)
	} else if !ta.RedirectUris()[rawRediUri] {
		return nil, erro.New("redirect URI " + rawRediUri + " is not registered")
	} else if uri, err := url.Parse(rawRediUri); err != nil {
		return nil, erro.Wrap(err)
	} else if stat := vals.Get(tagState); stat == "" {
		return uri, nil
	} else {
		q := uri.Query()
		q.Set(tagState, stat)
		uri.RawQuery = q.Encode()
		return uri, nil
	}
}
