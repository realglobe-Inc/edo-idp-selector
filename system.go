package main

import (
	"regexp"
	"strings"
)

type system struct {
	uiUri  string
	uiPath string

	cookMaxAge int

	idpCont idpContainer
}

func newSystem(uiUri string, uiPath string, cookMaxAge int, idpCont idpContainer) *system {
	uiUri = strings.TrimRight(uiUri, "/")
	uiUri = regexp.MustCompile("/+").ReplaceAllString(uiUri, "/")
	if uiUri == "" {
		uiUri = "/html"
	}
	if uiUri[0] != '/' {
		uiUri = "/" + uiUri
	}
	log.Info("Use " + uiUri + " as UI uri")
	return &system{uiUri, uiPath, cookMaxAge, idpCont}
}

func (sys *system) close() error {
	return sys.idpCont.close()
}
