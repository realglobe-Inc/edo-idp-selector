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
	"github.com/realglobe-Inc/edo-lib/server"
	"github.com/realglobe-Inc/go-lib/erro"
	"net/http"
	"net/url"
)

const (
	cookIdpId = "X-Edo-Idp-Id"
)

const (
	formIdpId   = "idp"
	formPrompt  = "prompt"
	formRediUri = "redirect_uri"
	formErr     = "error"
	formErrDesc = "error_description"
)

const (
	promptSelect = "select_account"
)

const (
	errInvalidReq = "invalid_request"
	errServerErr  = "server_error"
)

// redirect_uri があれば、リダイレクトしてエラーを通知する。
func handleError(w http.ResponseWriter, r *http.Request, err error) error {
	rediUri := r.FormValue(formRediUri)
	if rediUri == "" {
		// redirect_uri がない。
		return erro.Wrap(err)
	}

	// redirect_uri があった。

	log.Debug("Report error via redirect")

	redi, err := url.Parse(rediUri)
	if err != nil {
		return erro.Wrap(server.NewError(http.StatusBadRequest, "invalid "+formRediUri, err))
	}

	var errStr string
	switch e := erro.Unwrap(err).(type) {
	case *server.Error:
		switch e.Status() {
		case http.StatusBadRequest:
			errStr = errInvalidReq
		default:
			errStr = errServerErr
		}
	default:
		errStr = errServerErr
	}

	f := redi.Query()
	f.Set(formErr, errStr)
	f.Set(formErrDesc, erro.Unwrap(err).Error())
	redi.RawQuery = f.Encode()
	http.Redirect(w, r, redi.String(), http.StatusFound)
	return nil
}

// UI にリダイレクトする。
func redirectUi(sys *system, w http.ResponseWriter, r *http.Request, idp *idProvider) error {
	if idp != nil {
		// 補助としてデフォルト IdP を渡す。
		log.Debug("Default IdP is added to UI redirect uri")
		r.Form.Add(formIdpId, idp.Id)
	}

	query := r.Form.Encode()
	if query != "" {
		query = "?" + query
	}
	http.Redirect(w, r, sys.uiUri+query, http.StatusFound)
	return nil
}

// IdP にリダイレクトする。
func redirectIdp(sys *system, w http.ResponseWriter, r *http.Request, idp *idProvider) error {
	// cookie に記録しておく。
	cook := &http.Cookie{
		Name:   cookIdpId,
		Value:  idp.Id,
		MaxAge: sys.cookMaxAge,
	}
	http.SetCookie(w, cook)

	query := r.Form.Encode()
	if query != "" {
		query = "?" + query
	}
	http.Redirect(w, r, idp.AuthUri+query, http.StatusFound)
	return nil
}

// IdP 指定を読む。
func parseIdp(sys *system, r *http.Request) (*idProvider, error) {
	// クエリ優先。
	idpId := r.FormValue(formIdpId)
	if idpId != "" {
		log.Debug("IdP " + idpId + " is specified by form " + formIdpId)
		r.Form.Del(formIdpId)
	} else {
		cook, err := r.Cookie(cookIdpId)
		if err != nil {
			if err != http.ErrNoCookie {
				err = erro.Wrap(err)
				log.Err(erro.Unwrap(err))
				log.Debug(err)
			}
			return nil, nil
		} else {
			idpId = cook.Value
			log.Debug("IdP " + idpId + " is specified by cookie " + cookIdpId)
		}
	}

	idp, err := sys.idpCont.get(idpId)
	if err != nil {
		return nil, erro.Wrap(err)
	} else if idp == nil {
		log.Warn("Specified IdP " + idpId + " is not exist")
		return nil, nil
	}
	return idp, nil
}

// IdP 選択処理。
func selectPage(sys *system, w http.ResponseWriter, r *http.Request) error {
	idp, err := parseIdp(sys, r)
	if err != nil {
		return handleError(w, r, erro.Wrap(err))
	}

	if prompt := r.FormValue(formPrompt); prompt == promptSelect {
		// 強制アカウント選択なので UI に飛ばす。
		log.Debug("Redirect to UI because " + formPrompt + "=" + promptSelect)
		return redirectUi(sys, w, r, idp)
	} else if idp == nil {
		// 有効な IdP が指定されていないので UI に飛ばす。
		log.Debug("Redirect to UI because no valid IdP is specified")
		return redirectUi(sys, w, r, nil)
	}

	// 有効な IdP が指定されてた。
	log.Debug("Valid IdP " + idp.Id + " is specified")
	return redirectIdp(sys, w, r, idp)
}

// IdP 選択後であることを前提としたリダイレクト処理。
func redirectPage(sys *system, w http.ResponseWriter, r *http.Request) error {
	idp, err := parseIdp(sys, r)
	if err != nil {
		return handleError(w, r, erro.Wrap(err))
	}

	if idp == nil {
		// 有効な IdP が選択されていない。
		log.Debug("No valid IdP is specified")
		return handleError(w, r, erro.Wrap(server.NewError(http.StatusBadRequest, "no valid idp", nil)))
	}

	// 有効な IdP が指定されてた。
	log.Debug("Valid IdP " + idp.Id + " is specified")
	return redirectIdp(sys, w, r, idp)
}
