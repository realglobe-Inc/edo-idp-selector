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
	"github.com/realglobe-Inc/edo-idp-selector/request"
	"github.com/realglobe-Inc/go-lib/erro"
	"net/http"
	"net/url"
)

type startRequest struct {
	query_    string
	selForced bool
	disp      string
	langs     []string
}

func parseStartRequest(r *http.Request) (*startRequest, error) {
	query := r.URL.RawQuery
	if query == "" {
		return nil, erro.New("no parameters")
	}
	vals, err := url.ParseQuery(query)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return &startRequest{
		query_:    query,
		selForced: request.FormValueSet(vals.Get(tagPrompt))[tagSelect_account],
		disp:      vals.Get(tagDisplay),
		langs:     request.FormValues(vals.Get(tagUi_locales)),
	}, nil
}

func (this *startRequest) query() string {
	return this.query_
}

func (this *startRequest) selectForced() bool {
	return this.selForced
}

func (this *startRequest) display() string {
	return this.disp
}

func (this *startRequest) languages() []string {
	return this.langs
}
