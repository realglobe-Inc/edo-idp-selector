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
	"github.com/realglobe-Inc/edo-lib/server"
	"github.com/realglobe-Inc/go-lib/erro"
	"net/http"
)

// TA 情報を返す API。

const (
	tagClient_name = "client_name"
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
	req, err := newRequest(r, hndl.uriPrefix)
	if err != nil {
		return idperr.New(idperr.Invalid_request, "invalid request format", http.StatusBadRequest, erro.Wrap(err))
	}

	if req.ta() == "" {
		return idperr.New(idperr.Invalid_request, "no TA", http.StatusBadRequest, nil)
	}

	ta, err := hndl.taDb.Get(req.ta())
	if err != nil {
		return erro.Wrap(err)
	} else if ta == nil {
		return idperr.New(idperr.Invalid_request, "TA "+req.ta()+" is not found", http.StatusNotFound, nil)
	}

	// 提供する情報を選別する。
	info := map[string]interface{}{}
	for lang, name := range ta.Names() {
		tag := tagClient_name
		if lang != "" {
			tag += "#" + lang
		}
		info[tag] = name
	}

	return response(w, info)
}

// レスポンスを返す。
func response(w http.ResponseWriter, params map[string]interface{}) error {
	data, err := json.Marshal(params)
	if err != nil {
		return erro.Wrap(err)
	}

	w.Header().Add("Content-Type", server.ContentTypeJson)
	w.Write(data)
	return nil
}
