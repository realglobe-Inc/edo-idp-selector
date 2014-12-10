package main

import (
	"github.com/realglobe-Inc/edo/driver"
	"reflect"
	"testing"
	"time"
)

const (
	testIdpName = "test-id-provider-no-name"
)

// 要事前登録。
var testIdps = []*IdProvider{
	&IdProvider{Uuid: testIdpUuid + "-1", Name: testIdpName + "-1"},
	&IdProvider{Uuid: testIdpUuid + "-2", Name: testIdpName + "-2"},
}

func testIdpLister(t *testing.T, reg IdpLister) {
	idps, _, err := reg.IdProviders(nil)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(idps, testIdps) {
		t.Error(idps)
	}
}

func testIdpListerStamp(t *testing.T, reg IdpLister) {

	idps1, stmp1, err := reg.IdProviders(nil)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(idps1, testIdps) || stmp1 == nil {
		t.Error(idps1, stmp1)
	}

	// キャッシュと同じだから返らない。
	idps2, stmp2, err := reg.IdProviders(stmp1)
	if err != nil {
		t.Fatal(err)
	} else if idps2 != nil || stmp2 == nil {
		t.Error(idps2, stmp2)
	}

	// キャッシュが古いから返る。
	idps3, stmp3, err := reg.IdProviders(&driver.Stamp{Date: stmp1.Date.Add(-time.Second), Digest: stmp1.Digest})
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(idps3, testIdps) || stmp3 == nil {
		t.Error(idps3, stmp3)
	}

	// ダイジェストが違うから返る。
	idps4, stmp4, err := reg.IdProviders(&driver.Stamp{Date: stmp1.Date, Digest: stmp1.Digest + "a"})
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(idps4, testIdps) || stmp4 == nil {
		t.Error(idps4, stmp4)
	}
}
