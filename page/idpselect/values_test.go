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

package idpselect

import (
	idpdb "github.com/realglobe-Inc/edo-idp-selector/database/idp"
	tadb "github.com/realglobe-Inc/edo-idp-selector/database/ta"
	"github.com/realglobe-Inc/edo-lib/strset/strsetutil"
	"net/url"
)

const (
	test_label = "edo-test"
)

const (
	test_pathSelUi = "/ui/select.html"
	test_sessId    = "pbqy9Fx6OKqyGFe6FYS8SsqzZNuWxL"
	test_ticId     = "2IHV7qG7SJ"
	test_lang      = "ja"

	test_rediUri = "https://ta.example.org/callback"
)

var (
	test_idp1 = idpdb.New(
		"https://idp1.exampl.org",
		map[string]string{
			"":   "ID Provider 1",
			"ja": "ID プロバイダ 1号",
		},
		"https://idp1.exampl.org/auth",
		"", "", "", "", nil,
	)
	test_idp2 = idpdb.New(
		"https://idp2.exampl.org",
		map[string]string{
			"":   "ID Provider 2",
			"ja": "ID プロバイダ 2号",
		},
		"https://idp2.exampl.org/auth",
		"", "", "", "", nil,
	)
	test_idp3 = idpdb.New(
		"https://idp3.exampl.org",
		map[string]string{
			"":   "ID Provider 3",
			"ja": "ID プロバイダ 3号",
		},
		"https://idp3.exampl.org/auth",
		"", "", "", "", nil,
	)
	test_ta    = tadb.New("https://ta.example.org", nil, strsetutil.New(test_rediUri), nil, false, "")
	test_query = "response_type=code&scope=openid" +
		"&client_id=" + url.QueryEscape(test_ta.Id()) +
		"&redirect_uri=" + url.QueryEscape(test_rediUri)
)
