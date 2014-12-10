package main

import (
	"encoding/json"
	"github.com/realglobe-Inc/edo/driver"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"reflect"
	"testing"
	"time"
)

const (
	testLabel = "edo-test"

	testIdpUuid = "test-id-provider-no-uuid"

	testAttrName = "test-attribute-no-name"
)

var testAttr = map[string]interface{}{"array": []interface{}{"elem-1", "elem-2"}}

// JSON を通して等しいかどうか調べる。
func jsonEqual(v1 interface{}, v2 interface{}) (equal bool) {
	b1, err := json.Marshal(v1)
	if err != nil {
		log.Err(erro.Wrap(err))
		return false
	}
	var w1 interface{}
	if err := json.Unmarshal(b1, &w1); err != nil {
		log.Err(erro.Wrap(err))
		return false
	}

	b2, err := json.Marshal(v2)
	if err != nil {
		log.Err(erro.Wrap(err))
		return false
	}
	var w2 interface{}
	if err := json.Unmarshal(b2, &w2); err != nil {
		log.Err(erro.Wrap(err))
		return false
	}

	return reflect.DeepEqual(w1, w2)
}

// 要事前登録。

func testIdpAttributeProvider(t *testing.T, reg IdpAttributeProvider) {
	idpUuid := testIdpUuid
	attrName := testAttrName
	idpAttr := testAttr

	idpAttr1, _, err := reg.IdProviderAttribute(idpUuid, attrName, nil)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(idpAttr1, idpAttr) {
		if !jsonEqual(idpAttr1, idpAttr) {
			t.Error(idpAttr1)
		}
	}

	idpAttr2, _, err := reg.IdProviderAttribute(idpUuid, attrName+"1", nil)
	if err != nil {
		t.Fatal(err)
	} else if idpAttr2 != nil {
		t.Error(idpAttr2)
	}

	idpAttr3, _, err := reg.IdProviderAttribute(idpUuid+"1", attrName, nil)
	if err != nil {
		t.Fatal(err)
	} else if idpAttr3 != nil {
		t.Error(idpAttr3)
	}
}

func testIdpAttributeProviderStamp(t *testing.T, reg IdpAttributeProvider) {
	idpUuid := testIdpUuid
	attrName := testAttrName
	idpAttr := testAttr

	idpAttr1, stmp1, err := reg.IdProviderAttribute(idpUuid, attrName, nil)
	if err != nil {
		t.Fatal(err)
	} else if stmp1 == nil {
		t.Error(stmp1)
	} else if !reflect.DeepEqual(idpAttr1, idpAttr) {
		if !jsonEqual(idpAttr1, idpAttr) {
			t.Error(idpAttr1)
		}
	}

	idpAttr2, stmp2, err := reg.IdProviderAttribute(idpUuid, attrName+"1", nil)
	if err != nil {
		t.Fatal(err)
	} else if idpAttr2 != nil || stmp2 != nil {
		t.Error(idpAttr2, stmp2)
	}

	idpAttr3, stmp3, err := reg.IdProviderAttribute(idpUuid+"1", attrName, nil)
	if err != nil {
		t.Fatal(err)
	} else if idpAttr3 != nil || stmp3 != nil {
		t.Error(idpAttr3, stmp3)
	}

	// キャッシュと同じだから返らない。
	idpAttr4, stmp4, err := reg.IdProviderAttribute(idpUuid, attrName, stmp1)
	if err != nil {
		t.Fatal(err)
	} else if idpAttr4 != nil || stmp4 == nil {
		t.Error(idpAttr4, stmp4)
	}

	// キャッシュが古いから返る。
	idpAttr5, stmp5, err := reg.IdProviderAttribute(idpUuid, attrName, &driver.Stamp{Date: stmp1.Date.Add(-time.Second), Digest: stmp1.Digest})
	if err != nil {
		t.Fatal(err)
	} else if stmp5 == nil {
		t.Error(stmp5)
	} else if !reflect.DeepEqual(idpAttr5, idpAttr) {
		if !jsonEqual(idpAttr5, idpAttr) {
			t.Error(idpAttr5)
		}
	}

	// ダイジェストが違うから返る。
	idpAttr6, stmp6, err := reg.IdProviderAttribute(idpUuid, attrName, &driver.Stamp{Date: stmp1.Date, Digest: stmp1.Digest + "a"})
	if err != nil {
		t.Fatal(err)
	} else if stmp6 == nil {
		t.Error(stmp6)
	} else if !reflect.DeepEqual(idpAttr6, idpAttr) {
		if !jsonEqual(idpAttr6, idpAttr) {
			t.Error(idpAttr6)
		}
	}
}
