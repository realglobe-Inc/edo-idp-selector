package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileIdpLister(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	reg := NewFileIdpLister(path, 0)
	if _, err := reg.(*idpLister).base.Put("list", testIdps); err != nil {
		t.Fatal(err)
	}

	testIdpLister(t, reg)
}

func TestFileIdpListerStamp(t *testing.T) {
	path, err := ioutil.TempDir("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	reg := NewFileIdpLister(path, 0)
	if _, err := reg.(*idpLister).base.Put("list", testIdps); err != nil {
		t.Fatal(err)
	}

	testIdpListerStamp(t, reg)
}
