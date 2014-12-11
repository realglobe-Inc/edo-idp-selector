package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const filePerm = 0644

func TestFileIdpContainer(t *testing.T) {
	path, err := ioutil.TempDir("", "edo-idp-selector")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	idpCont := newFileIdpContainer(path, 0, 0)
	for _, idp := range []*idProvider{testIdp, testIdp2} {
		idpPath := filepath.Join(path, keyToEscapedJsonPath(idp.Id))
		buff, err := json.Marshal(idp)
		if err != nil {
			t.Fatal(err)
		}
		if err := ioutil.WriteFile(idpPath, buff, filePerm); err != nil {
			t.Fatal(err)
		}
	}
	testIdpContainer(t, idpCont)
}
