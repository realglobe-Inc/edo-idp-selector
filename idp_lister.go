package main

import (
	"github.com/realglobe-Inc/edo/driver"
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

// ID プロバイダ選択時に列挙する情報。
type IdProvider struct {
	Name string `json:"name" bson:"name"`
	Uuid string `json:"uuid" bson:"uuid"`
}

func (idp *IdProvider) String() string {
	return idp.Uuid + "," + idp.Name
}

type IdpLister interface {
	// ID プロバイダの列挙。
	IdProviders(caStmp *driver.Stamp) (idps []*IdProvider, newCaStmp *driver.Stamp, err error)
}

// 骨組み。
// バックエンドに ID プロバイダのリストそのものを保存。
type idpLister struct {
	base driver.KeyValueStore
}

func newIdpLister(base driver.KeyValueStore) *idpLister {
	return &idpLister{base: base}
}

func (reg *idpLister) IdProviders(caStmp *driver.Stamp) (idps []*IdProvider, newCaStmp *driver.Stamp, err error) {
	value, newCaStmp, err := reg.base.Get("list", caStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if value == nil {
		return nil, newCaStmp, nil
	}
	return value.([]*IdProvider), newCaStmp, nil
}
