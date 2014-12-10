package main

import (
	"github.com/realglobe-Inc/edo/driver"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2"
	"time"
)

// スレッドセーフ。
func NewMongoIdpLister(url, dbName, collName string, expiDur time.Duration) (IdpLister, error) {
	base, err := driver.NewMongoKeyValueStore(url, dbName, collName, expiDur)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	base.SetTake(func(query *mgo.Query) (interface{}, *driver.Stamp, error) {
		var res struct {
			Value []*IdProvider
			Stamp *driver.Stamp
		}
		if err := query.One(&res); err != nil {
			return nil, nil, erro.Wrap(err)
		}
		return res.Value, res.Stamp, nil
	})
	return newIdpLister(base), nil
}
