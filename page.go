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
	"github.com/realglobe-Inc/edo-idp-selector/database/session"
	idperr "github.com/realglobe-Inc/edo-idp-selector/error"
	"github.com/realglobe-Inc/edo-idp-selector/request"
	"github.com/realglobe-Inc/go-lib/erro"
	"net/http"
	"net/url"
)

func (sys *system) respondPageError(w http.ResponseWriter, r *http.Request, origErr error, sender *request.Request, sess *session.Element) (err error) {
	var uri *url.URL
	if sess.Query() != "" {
		uri, err = sys.getRedirectUri(sess.Query())
		if err != nil {
			log.Err(sender, ": ", erro.Unwrap(err))
			log.Debug(sender, ": ", erro.Wrap(err))
		}
	}

	// 経過を破棄。
	sess.Clear()
	if err := sys.sessDb.Save(sess, sess.Expires().Add(sys.sessDbExpIn-sys.sessExpIn)); err != nil {
		log.Err(sender, ": ", erro.Wrap(err))
	} else {
		log.Debug(sender, ": Saved session "+mosaic(sess.Id()))
	}

	if !sess.Saved() {
		// 未通知セッションの通知。
		http.SetCookie(w, sys.newCookie(sess))
		log.Debug(sender, ": Report session "+mosaic(sess.Id()))
	}

	if uri != nil {
		idperr.RedirectError(w, r, origErr, uri, sender)
	}

	return idperr.RespondPageError(w, r, origErr, sender, sys.errTmpl)
}

// リクエストからリダイレクト URI を取得する。
func (sys *system) getRedirectUri(rawQuery string) (*url.URL, error) {
	vals, err := url.ParseQuery(rawQuery)
	if err != nil {
		return nil, erro.Wrap(err)
	} else if taId := vals.Get(tagClient_id); taId == "" {
		return nil, erro.New("no TA ID")
	} else if rawRediUri := vals.Get(tagRedirect_uri); rawRediUri == "" {
		return nil, erro.New("no redirect URI")
	} else if ta, err := sys.taDb.Get(taId); err != nil {
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
