package main

import (
	"encoding/json"
	"github.com/realglobe-Inc/edo/driver"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"strings"
	"time"
)

func keyToJsonPath(key string) string {
	return key + ".json"
}

func jsonPathToKey(path string) string {
	if !strings.HasSuffix(path, ".json") {
		return ""
	}
	return path[:len(path)-len(".json")]
}

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
	return newIdpLister(driver.NewFileKeyValueStore(path, keyToJsonPath, jsonPathToKey, json.Marshal, idProvidersUnmarshal, expiDur, expiDur))
}
