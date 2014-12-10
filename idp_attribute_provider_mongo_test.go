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

func TestMongoIdpAttributeProvider(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoIdpAttributeProvider(mongoAddr, testLabel, "idp-attribute-provider", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*idpAttributeProvider).base.(driver.MongoKeyValueStore).Clear()

	if _, err := reg.(*idpAttributeProvider).base.Put(testIdpUuid+"/"+testAttrName, testAttr); err != nil {
		t.Fatal(err)
	}

	testIdpAttributeProvider(t, reg)
}

func TestMongoIdpAttributeProviderStamp(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoIdpAttributeProvider(mongoAddr, testLabel, "idp-attribute-provider", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*idpAttributeProvider).base.(driver.MongoKeyValueStore).Clear()

	if _, err := reg.(*idpAttributeProvider).base.Put(testIdpUuid+"/"+testAttrName, testAttr); err != nil {
		t.Fatal(err)
	}

	testIdpAttributeProviderStamp(t, reg)
}
