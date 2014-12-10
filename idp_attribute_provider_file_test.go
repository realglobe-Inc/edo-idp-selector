package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileIdpAttributeProvider(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	reg := NewFileIdpAttributeProvider(path, 0)
	if _, err := reg.(*idpAttributeProvider).base.Put(testIdpUuid+"/"+testAttrName, testAttr); err != nil {
		t.Fatal(err)
	}

	testIdpAttributeProvider(t, reg)
}

func TestFileIdpAttributeProviderStamp(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	reg := NewFileIdpAttributeProvider(path, 0)
	if _, err := reg.(*idpAttributeProvider).base.Put(testIdpUuid+"/"+testAttrName, testAttr); err != nil {
		t.Fatal(err)
	}

	testIdpAttributeProviderStamp(t, reg)
}
