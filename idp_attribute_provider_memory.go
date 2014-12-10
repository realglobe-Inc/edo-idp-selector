package main

import (
	"github.com/realglobe-Inc/edo/driver"
	"time"
)

type MemoryIdpAttributeProvider struct {
	base driver.KeyValueStore
}

// スレッドセーフ。
func NewMemoryIdpAttributeProvider(expiDur time.Duration) *MemoryIdpAttributeProvider {
	return &MemoryIdpAttributeProvider{driver.NewMemoryKeyValueStore(expiDur)}
}

func (reg *MemoryIdpAttributeProvider) IdProviderAttribute(idpUuid, attrName string, caStmp *driver.Stamp) (idpAttr interface{}, newCaStmp *driver.Stamp, err error) {
	return reg.base.Get(idpUuid+"/"+attrName, caStmp)
}

func (reg *MemoryIdpAttributeProvider) AddIdProviderAttribute(idpUuid, attrName string, idpAttr interface{}) {
	reg.base.Put(idpUuid+"/"+attrName, idpAttr)
}

func (reg *MemoryIdpAttributeProvider) RemoveIdProviderAttribute(idpUuid, attrName string) {
	reg.base.Remove(idpUuid + "/" + attrName)
}
