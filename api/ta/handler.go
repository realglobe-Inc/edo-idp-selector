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

// TA 情報を返す API。
package ta

import (
	"encoding/json"
	tadb "github.com/realglobe-Inc/edo-idp-selector/database/ta"
	idperr "github.com/realglobe-Inc/edo-idp-selector/error"
	"github.com/realglobe-Inc/edo-lib/server"
	"github.com/realglobe-Inc/go-lib/erro"
	"github.com/realglobe-Inc/go-lib/rglog/level"
	"net/http"
)

type handler struct {
	stopper *server.Stopper

	path string

	db tadb.Db

	debug bool
}

// path: 提供する URL パス。
// db: TA 情報 DB。
func New(
	stopper *server.Stopper,
	path string,
	db tadb.Db,
	debug bool,
) http.Handler {
	return &handler{
		stopper,
		path,
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

	//////////////////////////////
	server.LogRequest(level.DEBUG, r, this.debug)
	//////////////////////////////

	logPref = server.ParseSender(r) + ": "

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
	req, err := parseRequest(r, this.path)
	if err != nil {
		return erro.Wrap(idperr.New(idperr.Invalid_request, erro.Unwrap(err).Error(), http.StatusBadRequest, err))
	}

	log.Debug(this.logPref, "Parsed TA request")

	ta, err := this.db.Get(req.ta())
	if err != nil {
		return erro.Wrap(err)
	} else if ta == nil {
		return erro.Wrap(idperr.New(idperr.Invalid_request, "TA "+req.ta()+" is not exist", http.StatusNotFound, nil))
	}

	log.Debug(this.logPref, "Found TA "+req.ta())

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

	log.Debug(this.logPref, "Respond")

	if _, err := w.Write(data); err != nil {
		log.Err(this.logPref, erro.Wrap(err))
	}
	return nil
}
