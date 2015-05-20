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

import (
	idpdb "github.com/realglobe-Inc/edo-idp-selector/database/idp"
	"github.com/realglobe-Inc/edo-idp-selector/database/session"
	tadb "github.com/realglobe-Inc/edo-idp-selector/database/ta"
	webdb "github.com/realglobe-Inc/edo-idp-selector/database/web"
	"time"
)

func newTestSystem(webs []webdb.Element, idps []idpdb.Element, tas []tadb.Element) *system {
	return &system{
		pathSelUi: test_pathSelUi,

		sessLabel:    "Idp-Selector",
		sessLen:      30,
		sessExpIn:    time.Minute,
		sessRefDelay: time.Minute / 2,
		sessDbExpIn:  10 * time.Minute,
		ticLen:       10,
		ticExpIn:     time.Minute,

		webDb:  webdb.NewMemoryDb(webs),
		idpDb:  idpdb.NewMemoryDb(idps),
		taDb:   tadb.NewMemoryDb(tas),
		sessDb: session.NewMemoryDb(),

		cookPath: "/",
		cookSec:  false,
	}
}
