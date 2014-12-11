package main

import (
	"testing"
)

func TestMemoryIdpContainer(t *testing.T) {
	idpCont := newMemoryIdpContainer(0, 0)
	idpCont.add(testIdp)
	idpCont.add(testIdp2)
	testIdpContainer(t, idpCont)
}
