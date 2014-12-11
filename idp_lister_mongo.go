package main

import (
	"github.com/realglobe-Inc/edo/driver"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2"
	"time"
)

// スレッドセーフ。
func NewMongoIdpLister(url, dbName, collName string, expiDur time.Duration) (IdpLister, error) {
	return newIdpLister(driver.NewMongoKeyValueStore(url, dbName, collName, nil, nil, func(query *mgo.Query) (interface{}, *driver.Stamp, error) {
		var res struct {
			V []*IdProvider
			S *driver.Stamp
		}
		if err := query.One(&res); err != nil {
			return nil, nil, erro.Wrap(err)
		}
		return res.V, res.S, nil
	}, expiDur, expiDur)), nil
}
