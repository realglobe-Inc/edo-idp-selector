package main

import (
	"testing"
)

func TestMemoryIdpLister(t *testing.T) {
	reg := NewMemoryIdpLister(0)
	reg.SetIdProviders(testIdps)
	testIdpLister(t, reg)
}

func TestMemoryIdpListerStamp(t *testing.T) {
	reg := NewMemoryIdpLister(0)
	reg.SetIdProviders(testIdps)
	testIdpListerStamp(t, reg)
}
