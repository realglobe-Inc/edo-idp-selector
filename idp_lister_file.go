package main

import (
	"encoding/json"
	"github.com/realglobe-Inc/edo/driver"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// data を JSON として、[]*IdProvider にデコードする。
func idProvidersUnmarshal(data []byte) (interface{}, error) {
	var idps []*IdProvider
	if err := json.Unmarshal(data, &idps); err != nil {
		return nil, erro.Wrap(err)
	}
	return idps, nil
}

// スレッドセーフ。
func NewFileIdpLister(path string, expiDur time.Duration) IdpLister {
	return newIdpLister(driver.NewFileKeyValueStore(path, jsonKeyGen, json.Marshal, idProvidersUnmarshal, expiDur))
}
