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

package error

import (
	"github.com/realglobe-Inc/edo-idp-selector/request"
	"github.com/realglobe-Inc/edo-lib/server"
	"github.com/realglobe-Inc/go-lib/erro"
	"github.com/realglobe-Inc/go-lib/rglog/level"
	"html/template"
	"net/http"
)

func WrapPage(stopper *server.Stopper, f server.HandlerFunc, errTmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if stopper != nil {
			stopper.Stop()
			defer stopper.Unstop()
		}

		// panic 対策。
		defer func() {
			if rcv := recover(); rcv != nil {
				RespondHtml(w, r, erro.New(rcv), errTmpl, request.Parse(r, ""))
				return
			}
		}()

		//////////////////////////////
		server.LogRequest(level.DEBUG, r, true)
		//////////////////////////////

		if err := f(w, r); err != nil {
			RespondHtml(w, r, erro.Wrap(err), errTmpl, request.Parse(r, ""))
			return
		}
	}
}

func WrapApi(stopper *server.Stopper, f server.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if stopper != nil {
			stopper.Stop()
			defer stopper.Unstop()
		}

		// panic 対策。
		defer func() {
			if rcv := recover(); rcv != nil {
				RespondJson(w, r, erro.New(rcv), request.Parse(r, ""))
				return
			}
		}()

		//////////////////////////////
		server.LogRequest(level.DEBUG, r, true)
		//////////////////////////////

		if err := f(w, r); err != nil {
			RespondJson(w, r, erro.Wrap(err), request.Parse(r, ""))
			return
		}
	}
}