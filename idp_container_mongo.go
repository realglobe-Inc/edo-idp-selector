package main

import (
	"github.com/realglobe-Inc/edo/driver"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2"
	"time"
)

func readIdProvider(query *mgo.Query) (interface{}, *driver.Stamp, error) {
	var res struct {
		V *idProvider
		S *driver.Stamp
	}
	if err := query.One(&res); err != nil {
		return nil, nil, erro.Wrap(err)
	}
	return res.V, res.S, nil
}

// スレッドセーフ。
func newMongoIdpContainer(url, dbName, collName string, staleDur, expiDur time.Duration) idpContainer {
	return &idpContainerImpl{driver.NewMongoKeyValueStore(url, dbName, collName,
		nil, nil, readIdProvider,
		staleDur, expiDur)}
}
