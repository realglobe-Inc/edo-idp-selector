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
	"github.com/realglobe-Inc/edo-lib/jwk"
)

// ID プロバイダ情報。
type element struct {
	id       string
	names    map[string]string
	rediUris map[string]bool
	keys     []jwk.Key
	pw       bool
	sect     string
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
	return this.keys
}

func (this *element) Pairwise() bool {
	return this.pw
}

func (this *element) Sector() string {
	return this.sect
}
