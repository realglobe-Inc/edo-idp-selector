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
	"encoding/json"
	"github.com/realglobe-Inc/edo-idp-selector/database/session"
	idperr "github.com/realglobe-Inc/edo-idp-selector/error"
	"github.com/realglobe-Inc/edo-idp-selector/request"
	logutil "github.com/realglobe-Inc/edo-lib/log"
	"github.com/realglobe-Inc/go-lib/erro"
	"net/http"
	"net/url"
	"time"
)

func (this *Page) HandleStart(w http.ResponseWriter, r *http.Request) {
	var sender *request.Request

	// panic 対策。
	defer func() {
		if rcv := recover(); rcv != nil {
			idperr.RespondPageError(w, r, erro.New(rcv), sender, this.errTmpl)
			return
		}
	}()

	if this.stopper != nil {
		this.stopper.Stop()
		defer this.stopper.Unstop()
	}

	sender = request.Parse(r, this.sessLabel)
	log.Info(sender, ": Received start request")
	defer log.Info(sender, ": Handled start request")

	var sess *session.Element
	if sessId := sender.Session(); sessId != "" {
		// セッションが通知された。
		log.Debug(sender, ": Session is declared")

		var err error
		if sess, err = this.sessDb.Get(sessId); err != nil {
			log.Err(sender, ": ", erro.Wrap(err))
			// 新規発行すれば動くので諦めない。
		} else if sess == nil {
			// セッションが無かった。
			log.Warn(sender, ": Declared session is not exist")
		} else {
			// セッションがあった。
			log.Debug(sender, ": Declared session is exist")
		}
	}

	if now := time.Now(); sess == nil {
		// セッションを新規発行。
		sess = session.New(this.idGen.String(this.sessLen), now.Add(this.sessExpIn))
		log.Info(sender, ": Generated new session "+logutil.Mosaic(sess.Id())+" but not yet saved")
	} else if now.After(sess.Expires().Add(-this.sessRefDelay)) {
		// セッションを更新。
		old := sess
		sess = sess.New(this.idGen.String(this.sessLen), now.Add(this.sessExpIn))
		log.Info(sender, ": Refreshed session "+logutil.Mosaic(old.Id())+" to "+logutil.Mosaic(sess.Id())+" but not yet saved")
	}

	if err := this.startServe(w, r, sender, sess); err != nil {
		this.respondPageError(w, r, erro.Wrap(err), sender, sess)
		return
	}

	return
}

func (this *Page) startServe(w http.ResponseWriter, r *http.Request, sender *request.Request, sess *session.Element) error {
	req, err := parseStartRequest(r)
	if err != nil {
		return erro.Wrap(idperr.New(idperr.Invalid_request, erro.Unwrap(err).Error(), http.StatusBadRequest, err))
	}

	log.Debug(sender, ": Parsed start request")

	sess.SetQuery(req.query())

	if req.selectForced() {
		log.Debug(sender, ": Selection is forced")
	} else if sess.IdProvider() == "" {
		log.Debug(sender, ": Selection is required")
	} else {
		// 選択済み。
		if idp, err := this.idpDb.Get(sess.IdProvider()); err != nil {
			return erro.Wrap(err)
		} else if idp == nil {
			log.Warn(sender, ": Last selected ID provider "+sess.IdProvider()+" is not exist")
		} else {
			return this.redirectToIdProvider(w, r, idp, sender, sess)
		}
	}

	return this.redirectToSelectUi(w, r, req, sender, sess, "Please select your ID provider")
}

// 選択 UI にリダイレクトさせる。
func (this *Page) redirectToSelectUi(w http.ResponseWriter, r *http.Request, req *startRequest, sender *request.Request, sess *session.Element, msg string) error {
	uri, err := url.Parse(this.pathSelUi)
	if err != nil {
		return erro.Wrap(err)
	}

	// 選択 UI に渡すパラメータを生成。
	q := uri.Query()
	if idps := sess.SelectedIdProviders(); len(idps) > 0 {
		buff, err := json.Marshal(idps)
		if err != nil {
			return erro.Wrap(err)
		}
		q.Set(tagIssuers, string(buff))
	}
	if req.display() != "" {
		q.Set(tagDisplay, req.display())
	}
	if lang, langs := sess.Language(), req.languages(); lang != "" || len(langs) > 0 {
		a := []string{}
		m := map[string]bool{}
		for _, v := range append([]string{lang}, langs...) {
			if v == "" || m[v] {
				continue
			}
			a = append(a, v)
			m[v] = true
		}
		q.Set(tagLocales, request.ValuesForm(a))
	}
	if msg != "" {
		q.Set(tagMessage, msg)
	}
	uri.RawQuery = q.Encode()

	sess.SetTicket(session.NewTicket(this.idGen.String(this.ticLen), time.Now().Add(this.ticExpIn)))
	uri.Fragment = sess.Ticket().Id()
	log.Debug(sender, ": Published ticket "+logutil.Mosaic(sess.Ticket().Id()))

	if err := this.sessDb.Save(sess, sess.Expires().Add(this.sessDbExpIn-this.sessExpIn)); err != nil {
		log.Err(sender, ": ", erro.Wrap(err))
	} else {
		log.Debug(sender, ": Saved session "+logutil.Mosaic(sess.Id()))
	}

	if !sess.Saved() {
		http.SetCookie(w, this.newCookie(sess))
		log.Debug(sender, ": Report session "+logutil.Mosaic(sess.Id()))
	}

	log.Info(sender, ": Redirect to select UI")
	w.Header().Add(tagCache_control, tagNo_store)
	w.Header().Add(tagPragma, tagNo_cache)
	http.Redirect(w, r, uri.String(), http.StatusFound)
	return nil
}
