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

import ()

// コンパイル時に打ち間違いを検知するため。それ以上ではない。

const (
	// 小文字、アンダースコア。
	tagClient_id         = "client_id"
	tagDisplay           = "display"
	tagError             = "error"
	tagError_description = "error_description"
	tagIssuer            = "issuer"
	tagIssuer_name       = "issuer_name"
	tagIssuers           = "issuers"
	tagLocale            = "locale"
	tagLocales           = "locales"
	tagMessage           = "message"
	tagPrompt            = "prompt"
	tagRedirect_uri      = "redirect_uri"
	tagSelect_account    = "select_account"
	tagState             = "state"
	tagTicket            = "ticket"
	tagUi_locales        = "ui_locales"

	// 小文字、ハイフン。
	tagNo_cache = "no-cache"
	tagNo_store = "no-store"

	// 頭大文字、ハイフン。
	tagCache_control = "Cache-Control"
	tagContent_type  = "Content-Type"
	tagPragma        = "Pragma"
)

const (
	contTypeJson = "application/json"
)
