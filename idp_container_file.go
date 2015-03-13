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
	"encoding/json"
	"github.com/realglobe-Inc/edo-lib/driver"
	"github.com/realglobe-Inc/go-lib/erro"
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
	return &idpContainerImpl{driver.NewFileListedKeyValueStore(path,
		keyToEscapedJsonPath, escapedJsonPathToKey,
		json.Marshal, idProviderUnmarshal,
		staleDur, expiDur)}
}
