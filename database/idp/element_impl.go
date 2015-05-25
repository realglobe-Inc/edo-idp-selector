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
	"encoding/json"
	"github.com/realglobe-Inc/edo-idp-selector/database/web"
	"github.com/realglobe-Inc/edo-lib/jwk"
	"github.com/realglobe-Inc/go-lib/erro"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

// ID プロバイダ情報。
type element struct {
	id        string
	names     map[string]string
	authUri   string
	tokUri    string
	acntUri   string
	coopFrUri string
	coopToUri string
	keys      []jwk.Key

	// jwks_uri 用。
	keyUri string
	webDb  web.Db
}

// 主にテスト用。
func New(id string, names map[string]string, authUri, tokUri, acntUri, coopFrUri, coopToUri string, keys []jwk.Key) Element {
	return newElement(id, names, authUri, tokUri, acntUri, coopFrUri, coopToUri, keys)
}

func newElement(id string, names map[string]string, authUri, tokUri, acntUri, coopFrUri, coopToUri string, keys []jwk.Key) *element {
	return &element{
		id:        id,
		names:     names,
		authUri:   authUri,
		tokUri:    tokUri,
		acntUri:   acntUri,
		coopFrUri: coopFrUri,
		coopToUri: coopToUri,
		keys:      keys,
	}
}

func (this *element) Id() string {
	return this.id
}

func (this *element) Names() map[string]string {
	return this.names
}

func (this *element) AuthenticationUri() string {
	return this.authUri
}

func (this *element) TokenUri() string {
	return this.tokUri
}

func (this *element) AccountUri() string {
	return this.acntUri
}

func (this *element) CooperationFromUri() string {
	return this.coopFrUri
}

func (this *element) CooperationToUri() string {
	return this.coopToUri
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

// {
//     "issuer": <ID>,
//     "issuer_name": <表示名>,
//     "issuer_name#<言語タグ>": <表示名>,
//     ...,
//     "authorization_endpoint": <認証エンドポイント>,
//     "token_endpoint": <トークンエンドポイント>,
//     "userinfo_endpoint": <アカウント情報エンドポイント>,
//     "cooperation_from_endpoint": <要請元仲介エンドポイント>,
//     "cooperation_to_endpoint": <要請先仲介エンドポイント>,
//     "jwks": [
//         <JWK>,
//         ...
//     ],
//     "jwks_uri": <公開鍵の URI>
// }
func (this *element) SetBSON(raw bson.Raw) error {
	var buff struct {
		Id        string                   `bson:"issuer"`
		Names     map[string]interface{}   `bson:",inline"`
		AuthUri   string                   `bson:"authorization_endpoint"`
		TokUri    string                   `bson:"token_endpoint"`
		AcntUri   string                   `bson:"userinfo_endpoint"`
		CoopFrUri string                   `bson:"cooperation_from_endpoint"`
		CoopToUri string                   `bson:"cooperation_to_endpoint"`
		Keys      []map[string]interface{} `bson:"jwks"`
		KeyUri    string                   `bson:"jwks_uri"`
	}
	if err := raw.Unmarshal(&buff); err != nil {
		return erro.Wrap(err)
	}

	var names map[string]string
	for tag, name := range buff.Names {
		if !strings.HasPrefix(tag, "issuer_name") {
			continue
		}
		lang := tag[len("issuer_name"):]
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
	this.authUri = buff.AuthUri
	this.tokUri = buff.TokUri
	this.acntUri = buff.AcntUri
	this.coopFrUri = buff.CoopFrUri
	this.coopToUri = buff.CoopToUri
	this.keys = keys
	this.keyUri = buff.KeyUri
	return nil
}
