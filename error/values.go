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

package error

import ()

// error の値。
const (
	Access_denied              = "access_denied"
	Account_selection_required = "account_selection_required"
	Consent_required           = "consent_required"
	Interaction_required       = "interaction_required"
	Invalid_client             = "invalid_client"
	Invalid_grant              = "invalid_grant"
	Invalid_request            = "invalid_request"
	Invalid_scope              = "invalid_scope"
	Login_required             = "login_required"
	Registration_not_supported = "registration_not_supported"
	Request_not_supported      = "request_not_supported"
	Request_uri_not_supported  = "request_uri_not_supported"
	Server_error               = "server_error"
	Unsupported_grant_type     = "unsupported_grant_type"
	Unsupported_response_type  = "unsupported_response_type"
	// OpenID Connect の仕様ではサンプルとしてしか登場しない。
	Invalid_token = "invalid_token"
)

const (
	tagDebug             = "debug"
	tagError             = "error"
	tagError_description = "error_description"

	tagCache_control = "Cache-Control"
	tagContent_type  = "Content-Type"
	tagPragma        = "Pragma"

	tagNo_cache = "no-cache"
	tagNo_store = "no-store"
)

const (
	contTypeHtml = "text/html"
	contTypeJson = "application/json"
)
