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

// ID プロバイダ情報を返す API。
package idp

import (
	"encoding/json"
	"net/http"

	idpdb "github.com/realglobe-Inc/edo-idp-selector/database/idp"
	idperr "github.com/realglobe-Inc/edo-idp-selector/error"
	"github.com/realglobe-Inc/edo-lib/server"
	"github.com/realglobe-Inc/go-lib/erro"
	"github.com/realglobe-Inc/go-lib/rglog/level"
)

type handler struct {
	stopper *server.Stopper

	db idpdb.Db

	debug bool
}

func New(
	stopper *server.Stopper,
	db idpdb.Db,
	debug bool,
) http.Handler {
	return &handler{
		stopper,
		db,
		debug,
	}
}

func (this *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var logPref string

	// panic 対策。
	defer func() {
		if rcv := recover(); rcv != nil {
			idperr.RespondJson(w, r, erro.New(rcv), logPref)
			return
		}
	}()

	if this.stopper != nil {
		this.stopper.Stop()
		defer this.stopper.Unstop()
	}

	logPref = server.ParseSender(r) + ": "

	server.LogRequest(level.DEBUG, r, this.debug, logPref)

	log.Info(logPref, "Received TA request")
	defer log.Info(logPref, "Handled TA request")

	if err := (&environment{this, logPref}).serve(w, r); err != nil {
		idperr.RespondJson(w, r, erro.Wrap(err), logPref)
	}
}

// environment のメソッドは idperr.Error を返す。
type environment struct {
	*handler

	logPref string
}

func (this *environment) serve(w http.ResponseWriter, r *http.Request) error {
	req, err := parseRequest(r)
	if err != nil {
		return erro.Wrap(idperr.New(idperr.Invalid_request, erro.Unwrap(err).Error(), http.StatusBadRequest, err))
	}

	log.Debug(this.logPref, "Parsed ID provider request")

	idps, err := this.db.Search(req.filter())
	if err != nil {
		return erro.Wrap(err)
	}

	log.Debug(this.logPref, "Found ", len(idps), " ID providers")

	// 提供する情報を選別する。
	infos := []interface{}{}
	for _, idp := range idps {
		info := map[string]interface{}{
			tagIssuer: idp.Id(),
		}
		for lang, name := range idp.Names() {
			tag := tagIssuer_name
			if lang != "" {
				tag += "#" + lang
			}
			info[tag] = name
		}
		infos = append(infos, info)
	}

	data, err := json.Marshal(infos)
	if err != nil {
		return erro.Wrap(err)
	}

	w.Header().Add(tagContent_type, contTypeJson)

	log.Debug(this.logPref, "Repond")

	if _, err := w.Write(data); err != nil {
		log.Err(this.logPref, erro.Wrap(err))
	}
	return nil
}
