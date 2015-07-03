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
	"github.com/realglobe-Inc/edo-lib/server"
	"github.com/realglobe-Inc/go-lib/erro"
	"github.com/realglobe-Inc/go-lib/rglog/level"
	"net/http"
	"time"
)

func (this *Page) HandleSelect(w http.ResponseWriter, r *http.Request) {
	var logPref string

	// panic 対策。
	defer func() {
		if rcv := recover(); rcv != nil {
			idperr.RespondHtml(w, r, erro.New(rcv), this.errTmpl, logPref)
			return
		}
	}()

	if this.stopper != nil {
		this.stopper.Stop()
		defer this.stopper.Unstop()
	}

	//////////////////////////////
	server.LogRequest(level.DEBUG, r, this.debug)
	//////////////////////////////

	sender := request.Parse(r, this.sessLabel)
	logPref = sender.String() + ": "

	log.Info(logPref, "Received select request")
	defer log.Info(logPref, "Handled select request")

	var sess *session.Element
	if sessId := sender.Session(); sessId != "" {
		// セッションが通知された。
		log.Debug(logPref, "Session is declared")

		var err error
		if sess, err = this.sessDb.Get(sessId); err != nil {
			log.Err(logPref, erro.Wrap(err))
			// 新規発行すれば動くので諦めない。
		} else if sess == nil {
			// セッションが無かった。
			log.Warn(logPref, "Declared session is not exist")
		} else {
			// セッションがあった。
			log.Debug(logPref, "Declared session is exist")
		}
	}

	if sess == nil {
		// セッションを新規発行。
		sess = session.New(this.idGen.String(this.sessLen), time.Now().Add(this.sessExpIn))
		log.Info(logPref, "Generated new session "+logutil.Mosaic(sess.Id())+" but not yet saved")
	}

	env := (&environment{this, logPref, sess})
	if err := env.selectServe(w, r, sender); err != nil {
		env.respondErrorHtml(w, r, erro.Wrap(err))
		return
	}
}

func (this *environment) selectServe(w http.ResponseWriter, r *http.Request, sender *request.Request) error {
	req, err := parseSelectRequest(r)
	if err != nil {
		return erro.Wrap(idperr.New(idperr.Invalid_request, erro.Unwrap(err).Error(), http.StatusBadRequest, err))
	}

	log.Debug(this.logPref, "Parsed select request")

	if this.sess.Query() == "" {
		return erro.Wrap(idperr.New(idperr.Invalid_request, "not selecting session", http.StatusBadRequest, nil))
	} else if tic := this.sess.Ticket(); tic == nil {
		return erro.Wrap(idperr.New(idperr.Invalid_request, "ticket expired", http.StatusBadRequest, nil))
	} else if req.ticket() != tic.Id() {
		return erro.Wrap(idperr.New(idperr.Invalid_request, "invalid ticket", http.StatusBadRequest, nil))
	} else if tic.Expires().Before(time.Now()) {
		return erro.Wrap(idperr.New(idperr.Invalid_request, "ticket expired", http.StatusBadRequest, nil))
	}
	this.sess.SetTicket(nil)

	idp, err := this.idpDb.Get(req.idProvider())
	if err != nil {
		return erro.Wrap(err)
	} else if idp == nil {
		return erro.Wrap(idperr.New(idperr.Invalid_request, "ID provider "+req.idProvider()+" is not exist", http.StatusNotFound, nil))
	}

	this.sess.SelectIdProvider(idp.Id())
	log.Debug(this.logPref, "ID provider "+idp.Id()+" was selected")

	if lang := req.language(); lang != "" {
		this.sess.SetLanguage(lang)
		// 言語を選択してた。
		log.Debug(this.logPref, "Language "+lang+" was selected")
	}

	this.redirectToIdProvider(w, r, idp)
	return nil
}

func (this *environment) redirectToIdProvider(w http.ResponseWriter, r *http.Request, idp idpdb.Element) {
	uri := idp.AuthUri() + "?" + this.sess.Query()

	this.sess.Clear()
	if err := this.sessDb.Save(this.sess, this.sess.Expires().Add(this.sessDbExpIn-this.sessExpIn)); err != nil {
		log.Err(this.logPref, erro.Wrap(err))
	} else {
		log.Debug(this.logPref, "Saved session "+logutil.Mosaic(this.sess.Id()))
	}

	if !this.sess.Saved() {
		http.SetCookie(w, this.newCookie(this.sess))
		log.Debug(this.logPref, "Report session "+logutil.Mosaic(this.sess.Id()))
	}

	log.Info(this.logPref, "Redirect to ID provider "+idp.Id())
	w.Header().Add(tagCache_control, tagNo_store)
	w.Header().Add(tagPragma, tagNo_cache)
	http.Redirect(w, r, uri, http.StatusFound)
}
