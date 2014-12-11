package main

import (
	"github.com/realglobe-Inc/edo/driver"
	"time"
)

type MemoryIdpLister struct {
	base driver.KeyValueStore
}

// スレッドセーフ。
func NewMemoryIdpLister(expiDur time.Duration) *MemoryIdpLister {
	return &MemoryIdpLister{driver.NewMemoryKeyValueStore(expiDur, expiDur)}
}

func (reg *MemoryIdpLister) IdProviders(caStmp *driver.Stamp) (idps []*IdProvider, newCaStmp *driver.Stamp, err error) {
	value, newCaStmp, _ := reg.base.Get("list", caStmp)
	if value == nil {
		return nil, newCaStmp, nil
	}
	return value.([]*IdProvider), newCaStmp, nil
}

func (reg *MemoryIdpLister) SetIdProviders(idps []*IdProvider) {
	reg.base.Put("list", idps)
}
