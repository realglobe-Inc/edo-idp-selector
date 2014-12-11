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

func idProviderUnmarshal(data []byte) (interface{}, error) {
	var res idProvider
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, erro.Wrap(err)
	}
	return &res, nil
}

// スレッドセーフ。
func newFileIdpContainer(path string, staleDur, expiDur time.Duration) idpContainer {
	return &idpContainerImpl{driver.NewFileKeyValueStore(path,
		keyToEscapedJsonPath, escapedJsonPathToKey,
		json.Marshal, idProviderUnmarshal,
		staleDur, expiDur)}
}
