package main

import (
	"github.com/realglobe-Inc/edo/driver"
	"time"
)

// スレッドセーフ。
func NewMongoIdpAttributeProvider(url, dbName, collName string, expiDur time.Duration) (IdpAttributeProvider, error) {
	return newIdpAttributeProvider(driver.NewMongoKeyValueStore(url, dbName, collName, nil, nil, nil, expiDur, expiDur)), nil
}
