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
	"github.com/realglobe-Inc/edo-lib/jwk"
)

// ID プロバイダ情報。
type element struct {
	id        string
	names     map[string]string
	authUri   string
	coopFrUri string
	keys      []jwk.Key
}

func newElement(id string, names map[string]string, authUri, coopFrUri string, keys []jwk.Key) *element {
	return &element{
		id:        id,
		names:     names,
		authUri:   authUri,
		coopFrUri: coopFrUri,
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

func (this *element) CooperationFromUri() string {
	return this.coopFrUri
}

func (this *element) Keys() []jwk.Key {
	return this.keys
}
