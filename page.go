package main

import (
	"encoding/json"
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net/http"
	"net/url"
)

const (
	headerContentType = "Content-Type"
)

const cookieIdpId = "X-Edo-Idp-Id"

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
		return util.NewHttpStatusError(http.StatusBadRequest, "invalid "+formRediUri, erro.Wrap(err))
	}

	var errStr string
	switch e := erro.Unwrap(err).(type) {
	case *util.HttpStatusError:
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
	c := &http.Cookie{
		Name:   cookieIdpId,
		Value:  idp.Id,
		MaxAge: sys.cookieMaxAge,
	}
	http.SetCookie(w, c)

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
		c, err := r.Cookie(cookieIdpId)
		if err != nil {
			if err != http.ErrNoCookie {
				err = erro.Wrap(err)
				log.Err(erro.Unwrap(err))
				log.Debug(err)
			}
			return nil, nil
		} else {
			idpId = c.Value
			log.Debug("IdP " + idpId + " is specified by cookie " + cookieIdpId)
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

// IdP 一覧を返す。
func listApi(sys *system, w http.ResponseWriter, r *http.Request) error {
	// TODO クエリによる絞り込み。
	idps, err := sys.idpCont.list(nil)
	if err != nil {
		return erro.Wrap(err)
	}
	buff, err := json.Marshal(idps)
	if err != nil {
		return erro.Wrap(err)
	}
	log.Debug("Return ", len(idps), " IdPs")
	w.Header().Add(headerContentType, util.ContentTypeJson)
	w.Write(buff)
	return nil
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
		return handleError(w, r, erro.Wrap(util.NewHttpStatusError(http.StatusBadRequest, "no valid idp", nil)))
	}

	// 有効な IdP が指定されてた。
	log.Debug("Valid IdP " + idp.Id + " is specified")
	return redirectIdp(sys, w, r, idp)
}
