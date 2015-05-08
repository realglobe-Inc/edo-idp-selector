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
type Element interface {
	Id() string

	// 言語タグから表示名へのマップ。
	Names() map[string]string

	// 認証エンドポイント。
	AuthenticationUri() string

	// 要請元仲介エンドポイント。
	CooperationFromUri() string

	// 鍵。
	Keys() []jwk.Key
}
