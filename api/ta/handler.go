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

package ta

import (
	"encoding/json"
	tadb "github.com/realglobe-Inc/edo-idp-selector/database/ta"
	idperr "github.com/realglobe-Inc/edo-idp-selector/error"
	requtil "github.com/realglobe-Inc/edo-idp-selector/request"
	"github.com/realglobe-Inc/go-lib/erro"
	"net/http"
)

// TA 情報を返す API。

const (
	tagClient_name = "client_name"

	tagContent_type = "Content-Type"
)

const (
	contTypeJson = "application/json"
)

type Handler struct {
	// URI パスの接頭辞。
	uriPrefix string
	// TA 情報 DB。
	taDb tadb.Db
}

func NewHandler(uriPrefix string, taDb tadb.Db) *Handler {
	return &Handler{
		uriPrefix: uriPrefix,
		taDb:      taDb,
	}
}

func (hndl *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	sender := requtil.Parse(r, "")

	log.Info(sender, ": Received TA request")
	defer log.Info(sender, ": Handled TA request")

	if err := hndl.serve(w, r, sender); err != nil {
		return idperr.RespondApiError(w, r, erro.Wrap(err), sender)
	}
	return nil
}

func (hndl *Handler) serve(w http.ResponseWriter, r *http.Request, sender *requtil.Request) error {
	req, err := newRequest(r, hndl.uriPrefix)
	if err != nil {
		return erro.Wrap(idperr.New(idperr.Invalid_request, erro.Unwrap(err).Error(), http.StatusBadRequest, err))
	}

	log.Debug(sender, ": Parsed TA request")

	ta, err := hndl.taDb.Get(req.ta())
	if err != nil {
		return erro.Wrap(err)
	} else if ta == nil {
		return erro.Wrap(idperr.New(idperr.Invalid_request, "TA "+req.ta()+" is not exist", http.StatusNotFound, nil))
	}

	log.Debug(sender, ": Found TA "+req.ta())

	// 提供する情報を選別する。
	info := map[string]interface{}{}
	for lang, name := range ta.Names() {
		tag := tagClient_name
		if lang != "" {
			tag += "#" + lang
		}
		info[tag] = name
	}

	data, err := json.Marshal(info)
	if err != nil {
		return erro.Wrap(err)
	}

	w.Header().Add(tagContent_type, contTypeJson)

	log.Debug(sender, ": Return TA "+ta.Id())

	if _, err := w.Write(data); err != nil {
		log.Err(sender, ": ", erro.Wrap(err))
	}
	return nil
}
