package main

import (
	"github.com/realglobe-Inc/edo/driver"
)

type IdpAttributeProvider interface {
	// ID プロバイダの属性を返す。
	IdProviderAttribute(idpUuid, attrName string, caStmp *driver.Stamp) (idpAttr interface{}, newCaStmp *driver.Stamp, err error)
}

// 骨組み。
// バックエンドで ID プロバイダの属性ごとに保存。
type idpAttributeProvider struct {
	base driver.KeyValueStore
}

func newIdpAttributeProvider(base driver.KeyValueStore) *idpAttributeProvider {
	return &idpAttributeProvider{base}
}

func (reg *idpAttributeProvider) IdProviderAttribute(idpUuid, attrName string, caStmp *driver.Stamp) (idpAttr interface{}, newCaStmp *driver.Stamp, err error) {
	return reg.base.Get(idpUuid+"/"+attrName, caStmp)
}
