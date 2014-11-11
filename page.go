package main

import (
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net/http"
	"net/url"
)

const cookieIdpUuid = "ID_PROVIDER_UUID"

const (
	formRediUri     = "redirect_uri"
	formIdpLoginUri = "id_provider_login_uri"
	formIdpUuid     = "id_provider_uuid"
)

// /.
func routePage(sys *system, w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return erro.Wrap(err)
	}

	idpCookie, err := r.Cookie(cookieIdpUuid)
	if err != nil && err != http.ErrNoCookie {
		return erro.Wrap(err)
	}

	if idpCookie == nil || idpCookie.Value == "" {
		// 一覧ページに飛ばす。
		query := r.Form.Encode()
		if query != "" {
			query = "?" + query
		}
		w.Header().Set("Location", listPagePath+query)
		w.WriteHeader(http.StatusFound)
		log.Debug("No " + cookieIdpUuid + " in cookie.")
		return nil
	}

	// cookie に ID プロバイダが記録されてた。
	log.Debug(cookieIdpUuid + " " + idpCookie.Value + " in cookie.")

	idpLoginUri, _, err := sys.IdProviderLoginUri(idpCookie.Value, nil)
	if err != nil {
		return erro.Wrap(err)
	}

	if idpLoginUri == "" {
		// 一覧ページに飛ばす。
		query := r.Form.Encode()
		if query != "" {
			query = "?" + query
		}
		w.Header().Set("Location", listPagePath+query)
		w.WriteHeader(http.StatusFound)
		log.Debug("ID provider " + idpCookie.Value + " is invalid.")
		return nil
	}

	// 有効な ID プロバイダだった。
	log.Debug("ID provider " + idpCookie.Value + " is valid.")

	r.Form.Set(formIdpLoginUri, idpLoginUri)
	r.Form.Set(formIdpUuid, idpCookie.Value)
	query := r.Form.Encode()
	if query != "" {
		query = "?" + query
	}
	rediUri := setCookiePagePath + query
	w.Header().Set("Location", rediUri)
	w.WriteHeader(http.StatusFound)

	log.Debug("Redirect to " + rediUri + ".")
	return nil
}

// /list.
// ID プロバイダ一覧を表示する。
func listPage(sys *system, w http.ResponseWriter, r *http.Request) error {
	// TODO テンプレート HTML を読んで、そこに埋め込む形の方が良さそう。

	page := `
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.0 Transitional//EN">
<META HTTP-EQUIV="content-type" CONTENT="text/html; charset=utf-8">
<HTML>

  <HEAD>
    <TITLE>ID プロバイダ一覧</TITLE>
  </HEAD>

  <BODY>
    <H1>あなたの所属機関を選んでね。</H1>

    <P>`

	idps, _, err := sys.IdProviders(nil)
	if err != nil {
		return erro.Wrap(err)
	}

	if err := r.ParseForm(); err != nil {
		return erro.Wrap(err)
	}
	query := r.Form.Encode()
	if query != "" {
		query = "&" + query
	}

	for _, idp := range idps {
		form := url.Values{}
		form.Set(formIdpUuid, idp.Uuid)

		page += `
      <A HREF="` + setCookiePagePath + "?" + form.Encode() + query + "\">" + idp.Name + "</A></BR>"
	}

	page += `
    </P>
  </BODY>

</HTML>`

	w.Write([]byte(page))

	log.Debug("Responded list page.")
	return nil
}

// /set_cookie.
// cookie に ID プロバイダを記録してリダイレクトする。
func setCookiePage(sys *system, w http.ResponseWriter, r *http.Request) error {
	idpUuid := r.FormValue(formIdpUuid)
	if idpUuid == "" {
		return erro.Wrap(util.NewHttpStatusError(http.StatusBadRequest, "no "+formIdpUuid+" in parameters.", nil))
	}
	// id_provider_uuid は残す。

	// ID プロバイダの指定があった。
	log.Debug(formIdpUuid + " " + idpUuid + " in parameters.")

	idpLoginUri, _, err := sys.IdProviderLoginUri(idpUuid, nil)
	if err != nil {
		return erro.Wrap(err)
	} else if idpLoginUri == "" {
		return erro.Wrap(util.NewHttpStatusError(http.StatusForbidden, "ID provider "+idpUuid+" is invalid.", nil))
	}

	// 有効な ID プロバイダだった。
	log.Debug("ID provider " + idpUuid + " is valid.")

	idpCookie := &http.Cookie{
		Name:   cookieIdpUuid,
		Value:  idpUuid,
		MaxAge: sys.cookieMaxAge,
	}
	w.Header().Set("Set-Cookie", idpCookie.String())

	// リダイレクト先の URL クエリにも id_provider_uuid を付ける。
	var redi *url.URL
	if rediUri := r.FormValue(formRediUri); rediUri != "" {
		var err error
		redi, err = url.Parse(rediUri)
		if err != nil {
			return erro.Wrap(err)
		}

		f := redi.Query()
		f.Set(formIdpUuid, idpUuid)
		redi.RawQuery = f.Encode()
		r.Form.Set(formRediUri, redi.String())
	}

	query := r.Form.Encode()
	if query != "" {
		query = "?" + query
	}
	rediUri := idpLoginUri + query
	w.Header().Set("Location", rediUri)
	w.WriteHeader(http.StatusFound)

	log.Debug("Redirect to ID provider " + rediUri + ".")
	return nil
}
