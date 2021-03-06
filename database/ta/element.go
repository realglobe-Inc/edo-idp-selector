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

// TA 情報。
type Element interface {
	Id() string

	// 言語タグから表示名へのマップ。
	Names() map[string]string

	// リダイレクトエンドポイント。
	RedirectUris() map[string]bool

	// 鍵。
	Keys() []jwk.Key

	// セクタ固有のアカウント ID を使うかどうか。
	Pairwise() bool

	// セクタ固有のアカウント ID の計算に使うセクタ ID。
	Sector() string
}
