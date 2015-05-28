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
	"github.com/realglobe-Inc/go-lib/erro"
	"net/http"
)

type selectRequest struct {
	tic  string
	idp  string
	lang string
}

func parseSelectRequest(r *http.Request) (*selectRequest, error) {
	tic := r.FormValue(tagTicket)
	if tic == "" {
		return nil, erro.New("no ticket")
	}
	idp := r.FormValue(tagIssuer)
	if idp == "" {
		return nil, erro.New("no ID provider ID")
	}
	return &selectRequest{
		tic:  tic,
		idp:  idp,
		lang: r.FormValue(tagLocale),
	}, nil
}

func (this *selectRequest) ticket() string {
	return this.tic
}

func (this *selectRequest) idProvider() string {
	return this.idp
}

func (this *selectRequest) language() string {
	return this.lang
}
