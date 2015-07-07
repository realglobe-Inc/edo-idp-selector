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

package web

import (
	"io/ioutil"
	"net/http"

	"github.com/realglobe-Inc/go-lib/erro"
)

// web データを取ってくるだけ。
type directDb struct{}

func NewDirectDb() Db {
	return directDb{}
}

// 取得。
func (this directDb) Get(uri string) (Element, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, erro.New("cannot get data with status ", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	return newElement(uri, data), nil
}
