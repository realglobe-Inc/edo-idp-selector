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
	"html/template"
	"net/http"
	"time"
)

type system struct {
	pathSelUi string
	errTmpl   *template.Template

	sessLabel    string
	sessLen      int
	sessExpIn    time.Duration
	sessRefDelay time.Duration
	sessDbExpIn  time.Duration
	ticLen       int
	ticExpIn     time.Duration

	webDb  webdb.Db
	idpDb  idpdb.Db
	taDb   tadb.Db
	sessDb session.Db

	cookPath string
	cookSec  bool
}

func (sys *system) newCookie(sess *session.Element) *http.Cookie {
	return &http.Cookie{
		Name:     sys.sessLabel,
		Value:    sess.Id(),
		Path:     sys.cookPath,
		Expires:  sess.Expires(),
		Secure:   sys.cookSec,
		HttpOnly: true,
	}
}
