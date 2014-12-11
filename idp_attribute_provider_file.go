package main

import (
	"encoding/json"
	"github.com/realglobe-Inc/edo/driver"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net/url"
	"strings"
	"time"
)

func keyToEscapedJsonPath(key string) string {
	return url.QueryEscape(key) + ".json"
}

func escapedJsonPathToKey(path string) string {
	if !strings.HasSuffix(path, ".json") {
		return ""
	}
	key, _ := url.QueryUnescape(path[:len(path)-len(".json")])
	return key
}

// data を JSON として、encoding/json の標準データ型にデコードする。
func jsonUnmarshal(data []byte) (interface{}, error) {
	var res interface{}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, erro.Wrap(err)
	}
	return res, nil
}

// スレッドセーフ。
func NewFileIdpAttributeProvider(path string, expiDur time.Duration) IdpAttributeProvider {
	return newIdpAttributeProvider(driver.NewFileKeyValueStore(path, keyToEscapedJsonPath, escapedJsonPathToKey, json.Marshal, jsonUnmarshal, expiDur, expiDur))
}
