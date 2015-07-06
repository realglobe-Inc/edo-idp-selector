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
	"net/http"
	"net/url"
	"strings"

	"github.com/realglobe-Inc/go-lib/erro"
)

type request struct {
	ta_ string
}

func parseRequest(r *http.Request, uriPref string) (*request, error) {
	uriPref = strings.TrimRight(uriPref, "/") + "/"
	buff := strings.TrimPrefix(r.URL.Path, uriPref)
	if buff == "" {
		return nil, erro.New("no TA ID")
	}

	ta, err := url.QueryUnescape(buff)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	return &request{
		ta_: ta,
	}, nil
}

func (this *request) ta() string {
	return this.ta_
}
