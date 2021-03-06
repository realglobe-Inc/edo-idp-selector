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

package idp

import (
	"net/http"
	"net/url"

	"github.com/realglobe-Inc/go-lib/erro"
)

type request struct {
	filt map[string]string
}

func parseRequest(r *http.Request) (*request, error) {
	filt := map[string]string{}
	if r.URL.RawQuery != "" {
		vals, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			return nil, erro.Wrap(err)
		}
		for k, a := range vals {
			filt[k] = a[0]
		}
	}
	return &request{
		filt: filt,
	}, nil
}

func (this *request) filter() map[string]string {
	return this.filt
}
