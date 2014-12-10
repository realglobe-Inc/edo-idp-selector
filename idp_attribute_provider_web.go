package main

import (
	"encoding/json"
	"github.com/realglobe-Inc/edo/driver"
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

// {
//     "id_provider": {
//         attrNameX: XXX
//     }
// }
func webIdProviderAttributeUnmarshal(data []byte) (interface{}, error) {
	var res struct {
		Id_provider map[string]interface{}
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, erro.Wrap(err)
	}
	return res.Id_provider, nil
}

type webIdpAttributeProvider struct {
	driver.KeyValueStore
}

// スレッドセーフ。
func NewWebIdpAttributeProvider(prefix string) IdpAttributeProvider {
	return newWebIdpAttributeProvider(driver.NewWebKeyValueStore(prefix, nil, webIdProviderAttributeUnmarshal))
}

func newWebIdpAttributeProvider(base driver.KeyValueStore) *webIdpAttributeProvider {
	return &webIdpAttributeProvider{base}
}

func (reg *webIdpAttributeProvider) IdProviderAttribute(idpUuid, attrName string, caStmp *driver.Stamp) (idpAttr interface{}, newCaStmp *driver.Stamp, err error) {
	value, newCaStmp, err := reg.Get(idpUuid+"/"+attrName, caStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if value == nil {
		return nil, nil, nil
	}
	return value.(map[string]interface{})[attrName], newCaStmp, nil
}
