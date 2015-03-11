package main

import (
	"encoding/json"
	"github.com/realglobe-Inc/edo-lib/server"
	"github.com/realglobe-Inc/go-lib/erro"
	"net/http"
)

const (
	headContentType = "Content-Type"
)

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
	w.Header().Add(headContentType, server.ContentTypeJson)
	w.Write(buff)
	return nil
}
