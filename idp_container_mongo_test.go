package main

import (
	"github.com/realglobe-Inc/edo/driver"
	"gopkg.in/mgo.v2"
	"testing"
)

// テストするなら、ローカルにデフォルトポートで mongodb をたてる必要あり。
var mongoAddr = "localhost"

func init() {
	if mongoAddr != "" {
		// 実際にサーバーが立っているかどうか調べる。
		// 立ってなかったらテストはスキップ。
		conn, err := mgo.Dial(mongoAddr)
		if err != nil {
			mongoAddr = ""
		} else {
			conn.Close()
		}
	}
}

func TestMongoIdpContainer(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	idpCont := newMongoIdpContainer(mongoAddr, "edo-test", "edo-idp-selector", 0, 0)
	defer idpCont.(*idpContainerImpl).base.(driver.MongoKeyValueStore).Clear()
	for _, idp := range []*idProvider{testIdp, testIdp2} {
		if _, err := idpCont.(*idpContainerImpl).base.Put(idp.Id, idp); err != nil {
			t.Fatal(err)
		}
	}
	testIdpContainer(t, idpCont)
}
