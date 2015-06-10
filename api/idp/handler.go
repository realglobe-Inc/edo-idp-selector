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
	idpdb "github.com/realglobe-Inc/edo-idp-selector/database/idp"
	idperr "github.com/realglobe-Inc/edo-idp-selector/error"
	requtil "github.com/realglobe-Inc/edo-idp-selector/request"
	"github.com/realglobe-Inc/edo-lib/server"
	"github.com/realglobe-Inc/go-lib/erro"
	"net/http"
)

type handler struct {
	stopper *server.Stopper
	db      idpdb.Db
}

func New(stopper *server.Stopper, db idpdb.Db) http.Handler {
	return &handler{
		stopper: stopper,
		db:      db,
	}
}

func (this *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var sender *requtil.Request

	// panic 対策。
	defer func() {
		if rcv := recover(); rcv != nil {
			idperr.RespondJson(w, r, erro.New(rcv), sender)
			return
		}
	}()

	if this.stopper != nil {
		this.stopper.Stop()
		defer this.stopper.Unstop()
	}

	sender = requtil.Parse(r, "")
	log.Info(sender, ": Received TA request")
	defer log.Info(sender, ": Handled TA request")

	if err := this.serve(w, r, sender); err != nil {
		idperr.RespondJson(w, r, erro.Wrap(err), sender)
	}
}

func (this *handler) serve(w http.ResponseWriter, r *http.Request, sender *requtil.Request) error {
	req, err := parseRequest(r)
	if err != nil {
		return erro.Wrap(idperr.New(idperr.Invalid_request, erro.Unwrap(err).Error(), http.StatusBadRequest, err))
	}

	log.Debug(sender, ": Parsed ID provider request")

	idps, err := this.db.Search(req.filter())
	if err != nil {
		return erro.Wrap(err)
	}

	log.Debug(sender, ": Found ", len(idps), " ID providers")

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

	log.Debug(sender, ": Repond")

	if _, err := w.Write(data); err != nil {
		log.Err(sender, ": ", erro.Wrap(err))
	}
	return nil
}
