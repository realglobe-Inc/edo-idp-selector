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
	"github.com/realglobe-Inc/go-lib/erro"
	"net/http"
	"net/url"
	"strings"
)

type request struct {
	ta_ string
}

func newRequest(r *http.Request, uriPrefix string) (*request, error) {
	if len(uriPrefix) == 0 || uriPrefix[len(uriPrefix)-1] != '/' {
		uriPrefix += "/"
	}
	rawTa := strings.TrimPrefix(r.URL.Path, uriPrefix)
	ta, err := url.QueryUnescape(rawTa)
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