package main

import (
	"testing"
)

func TestMemoryIdpAttributeProvider(t *testing.T) {
	reg := NewMemoryIdpAttributeProvider(0)
	reg.AddIdProviderAttribute(testIdpUuid, testAttrName, testAttr)
	testIdpAttributeProvider(t, reg)
}

func TestMemoryIdpAttributeProviderStamp(t *testing.T) {
	reg := NewMemoryIdpAttributeProvider(0)
	reg.AddIdProviderAttribute(testIdpUuid, testAttrName, testAttr)
	testIdpAttributeProviderStamp(t, reg)
}
