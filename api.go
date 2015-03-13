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
	"github.com/realglobe-Inc/edo-lib/server"
	"github.com/realglobe-Inc/go-lib/erro"
	"net/http"
)

const (
	headContentType = "Content-Type"
)

// IdP 一覧を返す。
func listApi(sys *system, w http.ResponseWriter, r *http.Request) error {
	// TODO クエリによる絞り込み。
	idps, err := sys.idpCont.list(nil)
	if err != nil {
		return erro.Wrap(err)
	}
	buff, err := json.Marshal(idps)
	if err != nil {
		return erro.Wrap(err)
	}
	log.Debug("Return ", len(idps), " IdPs")
	w.Header().Add(headContentType, server.ContentTypeJson)
	w.Write(buff)
	return nil
}
