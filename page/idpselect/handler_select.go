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
	idperr "github.com/realglobe-Inc/edo-idp-selector/error"
	"github.com/realglobe-Inc/edo-idp-selector/request"
	logutil "github.com/realglobe-Inc/edo-lib/log"
	"github.com/realglobe-Inc/go-lib/erro"
	"net/http"
	"time"
)

func (this *Page) HandleSelect(w http.ResponseWriter, r *http.Request) {
	var sender *request.Request

	// panic 対策。
	defer func() {
		if rcv := recover(); rcv != nil {
			idperr.RespondHtml(w, r, erro.New(rcv), sender, this.errTmpl)
			return
		}
	}()

	if this.stopper != nil {
		this.stopper.Stop()
		defer this.stopper.Unstop()
	}

	sender = request.Parse(r, this.sessLabel)
	log.Info(sender, ": Received select request")
	defer log.Info(sender, ": Handled select request")

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

	if sess == nil {
		// セッションを新規発行。
		sess = session.New(this.idGen.String(this.sessLen), time.Now().Add(this.sessExpIn))
		log.Info(sender, ": Generated new session "+logutil.Mosaic(sess.Id())+" but not yet saved")
	}

	if err := this.selectServe(w, r, sender, sess); err != nil {
		this.respondPageError(w, r, erro.Wrap(err), sender, sess)
		return
	}

	return
}

func (this *Page) selectServe(w http.ResponseWriter, r *http.Request, sender *request.Request, sess *session.Element) error {
	req, err := parseSelectRequest(r)
	if err != nil {
		return erro.Wrap(idperr.New(idperr.Invalid_request, erro.Unwrap(err).Error(), http.StatusBadRequest, err))
	}

	log.Debug(sender, ": Parsed select request")

	if sess.Query() == "" {
		return erro.Wrap(idperr.New(idperr.Invalid_request, "not selecting session", http.StatusBadRequest, nil))
	} else if tic := sess.Ticket(); tic == nil {
		return erro.Wrap(idperr.New(idperr.Invalid_request, "ticket expired", http.StatusBadRequest, nil))
	} else if req.ticket() != tic.Id() {
		return erro.Wrap(idperr.New(idperr.Invalid_request, "invalid ticket", http.StatusBadRequest, nil))
	} else if tic.Expires().Before(time.Now()) {
		return erro.Wrap(idperr.New(idperr.Invalid_request, "ticket expired", http.StatusBadRequest, nil))
	}
	sess.SetTicket(nil)

	idp, err := this.idpDb.Get(req.idProvider())
	if err != nil {
		return erro.Wrap(err)
	} else if idp == nil {
		return erro.Wrap(idperr.New(idperr.Invalid_request, "ID provider "+req.idProvider()+" is not exist", http.StatusNotFound, nil))
	}

	sess.SelectIdProvider(idp.Id())
	log.Debug(sender, ": ID provider "+idp.Id()+" was selected")

	if lang := req.language(); lang != "" {
		sess.SetLanguage(lang)
		// 言語を選択してた。
		log.Debug(sender, ": Language "+lang+" was selected")
	}

	return this.redirectToIdProvider(w, r, idp, sender, sess)
}

func (this *Page) redirectToIdProvider(w http.ResponseWriter, r *http.Request, idp idpdb.Element, sender *request.Request, sess *session.Element) error {
	uri := idp.AuthUri() + "?" + sess.Query()

	sess.Clear()
	if err := this.sessDb.Save(sess, sess.Expires().Add(this.sessDbExpIn-this.sessExpIn)); err != nil {
		log.Err(sender, ": ", erro.Wrap(err))
	} else {
		log.Debug(sender, ": Saved session "+logutil.Mosaic(sess.Id()))
	}

	if !sess.Saved() {
		http.SetCookie(w, this.newCookie(sess))
		log.Debug(sender, ": Report session "+logutil.Mosaic(sess.Id()))
	}

	log.Info(sender, ": Redirect to ID provider "+idp.Id())
	w.Header().Add(tagCache_control, tagNo_store)
	w.Header().Add(tagPragma, tagNo_cache)
	http.Redirect(w, r, uri, http.StatusFound)
	return nil
}
