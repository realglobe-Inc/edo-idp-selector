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

import ()

// ID プロバイダ情報。
type Element interface {
	Id() string

	// 表示名。
	// 優先言語 lang が空なら任意の 1 つ。
	// 優先言語があっても、それが返るとは限らない。
	Name(langs []string) (name, lang string)

	// 認証エンドポイント。
	AuthUri() string

	// 要請元仲介エンドポイント。
	CoopSrcUri() string

	// 署名検証鍵。
	// 返り値は kid 値から鍵へのマップ。
	VerifyKeys() map[string]interface{}
}
