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
	"encoding/json"
	"github.com/realglobe-Inc/edo-idp-selector/database/web"
	"github.com/realglobe-Inc/edo-lib/jwk"
	"github.com/realglobe-Inc/edo-lib/strset"
	"github.com/realglobe-Inc/go-lib/erro"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

// ID プロバイダ情報。
type element struct {
	id       string
	names    map[string]string
	rediUris map[string]bool
	keys     []jwk.Key
	pw       bool
	sect     string

	// jwks_uri 用。
	keyUri string
	webDb  web.Db
}

// 主のテスト用。
func New(id string, names map[string]string, rediUris map[string]bool, keys []jwk.Key, pw bool, sect string) Element {
	return newElement(id, names, rediUris, keys, pw, sect)
}

func newElement(id string, names map[string]string, rediUris map[string]bool, keys []jwk.Key, pw bool, sect string) *element {
	return &element{
		id:       id,
		names:    names,
		rediUris: rediUris,
		keys:     keys,
		pw:       pw,
		sect:     sect,
	}
}

func (this *element) Id() string {
	return this.id
}

func (this *element) Names() map[string]string {
	return this.names
}

func (this *element) RedirectUris() map[string]bool {
	return this.rediUris
}

func (this *element) Keys() []jwk.Key {
	if this.keyUri != "" {
		if this.downloadKeys() {
			this.keyUri = ""
			this.webDb = nil
		}
	}
	return this.keys
}

func (this *element) Pairwise() bool {
	return this.pw
}

func (this *element) Sector() string {
	return this.sect
}

func (this *element) setWebDbIfNeeded(webDb web.Db) {
	if this.keyUri != "" {
		this.webDb = webDb
	}
}

func (this *element) downloadKeys() (ok bool) {
	if this.webDb == nil {
		return false
	}

	elem, err := this.webDb.Get(this.keyUri)
	if err != nil {
		log.Warn(erro.Wrap(err))
		return false
	} else if elem == nil {
		// そんなもの無かった。
		return true
	}

	var ma []map[string]interface{}
	if err := json.Unmarshal(elem.Data(), &ma); err != nil {
		log.Warn(erro.Wrap(err))
		return false
	}

	for _, m := range ma {
		key, err := jwk.FromMap(m)
		if err != nil {
			log.Warn(erro.Wrap(err))
			return false
		}
		this.keys = append(this.keys, key)
	}

	return true
}

//  {
//      "client_id": <ID>,
//      "client_name": <表示名>,
//      "client_name#<言語タグ>": <表示名>,
//      ...,
//      "redirect_uris": [
//          <リダイレクトエンドポイント>,
//          ...
//      ],
//      "jwks": [
//          <JWK>,
//          ...
//      ],
//      "jwks_uri": <公開鍵の URI>
//      "subject_type": <pairwise / public>,
//      "sector_identifier_uri": <セクタ ID>
// }
func (this *element) SetBSON(raw bson.Raw) error {
	var buff struct {
		Id       string                   `bson:"client_id"`
		Names    map[string]interface{}   `bson:",inline"`
		RediUris strset.Set               `bson:"redirect_uris"`
		Keys     []map[string]interface{} `bson:"jwks"`
		Pw       string                   `bson:"subject_type"`
		Sect     string                   `bson:"sector_identifier_uri"`
		KeyUri   string                   `bson:"jwks_uri"`
	}
	if err := raw.Unmarshal(&buff); err != nil {
		return erro.Wrap(err)
	}

	var names map[string]string
	for tag, name := range buff.Names {
		if !strings.HasPrefix(tag, "client_name") {
			continue
		}
		lang := tag[len("client_name"):]
		if len(lang) > 0 {
			if lang[0] != '#' {
				continue
			}
			lang = lang[1:]
		}
		if names == nil {
			names = map[string]string{}
		}
		names[lang], _ = name.(string)
	}
	var keys []jwk.Key
	if buff.Keys != nil {
		keys = []jwk.Key{}
		for _, m := range buff.Keys {
			key, err := jwk.FromMap(m)
			if err != nil {
				return erro.Wrap(err)
			}
			keys = append(keys, key)
		}
	}

	this.id = buff.Id
	this.names = names
	this.rediUris = buff.RediUris
	this.keys = keys
	this.pw = !(buff.Pw == "public")
	this.sect = buff.Sect
	this.keyUri = buff.KeyUri
	return nil
}
