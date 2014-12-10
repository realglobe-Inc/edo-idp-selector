package main

import (
	"github.com/realglobe-Inc/edo/driver"
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

const attrLoginUri = "login_uri"

// 便宜的に集めただけ。
type system struct {
	IdpLister
	IdpAttributeProvider

	cookieMaxAge int
}

func (sys *system) IdProviderLoginUri(idpUuid string, caStmp *driver.Stamp) (loginUri string, newCaStmp *driver.Stamp, err error) {
	value, newCaStmp, err := sys.IdpAttributeProvider.IdProviderAttribute(idpUuid, attrLoginUri, caStmp)
	if err != nil {
		return "", nil, erro.Wrap(err)
	} else if value == nil {
		return "", newCaStmp, nil
	}
	return value.(string), newCaStmp, nil
}
