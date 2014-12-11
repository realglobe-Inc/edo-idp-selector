package main

import (
	"reflect"
	"testing"
)

var testIdp = &idProvider{
	Id:      "https://example.com",
	Name:    "sample idp",
	AuthUri: "https://example.com/login",
}
var testIdp2 = &idProvider{
	Id:      "idp-no-id",
	Name:    "認証装置2",
	AuthUri: "https://a.b.c.example.com/",
}

func testIdpContainer(t *testing.T, idpCont idpContainer) {
	if idp, err := idpCont.get(testIdp.Id); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(idp, testIdp) {
		t.Error(idp)
	}

	if idps, err := idpCont.list(nil); err != nil {
		t.Fatal(err)
	} else if len(idps) != 2 {
		t.Error(idps)
	}
}
