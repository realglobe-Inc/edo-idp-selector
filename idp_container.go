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
	"github.com/realglobe-Inc/edo-lib/driver"
	"github.com/realglobe-Inc/go-lib/erro"
	"regexp"
)

type idProvider struct {
	Id      string `json:"id"                     bson:"id"`
	Name    string `json:"name"                   bson:"name"`
	AuthUri string `json:"authorization_endpoint" bson:"authorization_endpoint"`
}

type idpContainer interface {
	// 1 個取得。
	get(idpId string) (*idProvider, error)

	// まとめて取得。
	// filter はタグ名から該当する値の正規表現へのマップ。
	// filter の項目は「かつ」で結合。
	list(filter map[string]string) ([]*idProvider, error)

	close() error
}

type idpContainerImpl struct {
	base driver.ListedKeyValueStore
}

func (this *idpContainerImpl) get(idpId string) (*idProvider, error) {
	val, _, err := this.base.Get(idpId, nil)
	if err != nil {
		return nil, erro.Wrap(err)
	} else if val == nil {
		return nil, nil
	}
	return val.(*idProvider), nil
}

func (this *idpContainerImpl) list(filter map[string]string) ([]*idProvider, error) {
	cf := map[string]*regexp.Regexp{}
	for k, v := range filter {
		reg, err := regexp.Compile(v)
		if err != nil {
			return nil, erro.Wrap(err)
		}
		cf[k] = reg
	}

	keys, _, err := this.base.Keys(nil)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	idps := []*idProvider{}
	for key, _ := range keys {
		idp, err := this.get(key)
		if err != nil {
			return nil, erro.Wrap(err)
		} else if idp == nil {
			continue
		}

		buff := map[string]string{"id": idp.Id, "name": idp.Name, "authorization_endpoint": idp.AuthUri}
		ok := true
		for k, reg := range cf {
			if !reg.MatchString(buff[k]) {
				ok = false
				break
			}
		}
		if ok {
			idps = append(idps, idp)
		}
	}
	return idps, nil
}

func (this *idpContainerImpl) close() error {
	return this.base.Close()
}
