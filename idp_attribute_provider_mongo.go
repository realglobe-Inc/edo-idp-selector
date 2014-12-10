package main

import (
	"github.com/realglobe-Inc/edo/driver"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// スレッドセーフ。
func NewMongoIdpAttributeProvider(url, dbName, collName string, expiDur time.Duration) (IdpAttributeProvider, error) {
	base, err := driver.NewMongoKeyValueStore(url, dbName, collName, expiDur)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return newIdpAttributeProvider(base), nil
}
