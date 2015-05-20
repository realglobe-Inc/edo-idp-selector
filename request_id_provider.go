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
	"github.com/realglobe-Inc/go-lib/erro"
	"net/http"
	"net/url"
)

type idProviderRequest struct {
	filter_ map[string]string
}

func parseIdProviderRequest(r *http.Request) (*idProviderRequest, error) {
	filter := map[string]string{}
	if r.URL.RawQuery != "" {
		vals, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			return nil, erro.Wrap(err)
		}
		for k, a := range vals {
			filter[k] = a[0]
		}
	}
	return &idProviderRequest{
		filter_: filter,
	}, nil
}

func (this *idProviderRequest) filter() map[string]string {
	return this.filter_
}
