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
	"encoding/json"
	idperr "github.com/realglobe-Inc/edo-idp-selector/error"
	"github.com/realglobe-Inc/edo-idp-selector/request"
	"github.com/realglobe-Inc/go-lib/erro"
	"net/http"
)

// ID プロバイダ情報を返す API。

func (sys *system) idProviderApi(w http.ResponseWriter, r *http.Request) error {
	sender := request.Parse(r, "")

	log.Info(sender, ": Received ID provider request")
	defer log.Info(sender, ": Handled ID provider request")

	if err := sys.idProviderServe(w, r, sender); err != nil {
		return idperr.RespondApiError(w, r, erro.Wrap(err), sender)
	}
	return nil
}

func (sys *system) idProviderServe(w http.ResponseWriter, r *http.Request, sender *request.Request) error {
	req, err := parseIdProviderRequest(r)
	if err != nil {
		return erro.Wrap(idperr.New(idperr.Invalid_request, erro.Unwrap(err).Error(), http.StatusBadRequest, err))
	}

	log.Debug(sender, ": Parsed ID provider request")

	idps, err := sys.idpDb.Search(req.filter())
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

	log.Debug(sender, ": Return ", len(infos), " ID providers")

	if _, err := w.Write(data); err != nil {
		log.Err(sender, ": ", erro.Wrap(err))
	}
	return nil
}
