package main

import (
	"encoding/json"
	"github.com/realglobe-Inc/edo/driver"
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

// {
//     "id_providers": [
//         id-provider-no-hitotsu,
//         ...
//     ]
// }
func webIdProvidersUnmarshal(data []byte) (interface{}, error) {
	var res struct {
		Id_providers []*IdProvider
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, erro.Wrap(err)
	}
	return res.Id_providers, nil
}

// スレッドセーフ。
func NewWebIdpLister(prefix string) IdpLister {
	return newIdpLister(driver.NewWebKeyValueStore(prefix, nil, webIdProvidersUnmarshal))
}
