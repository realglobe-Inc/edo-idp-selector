package main

import (
	"regexp"
	"strings"
)

type system struct {
	uiUri  string
	uiPath string

	cookieMaxAge int

	idpCont idpContainer
}

var slashSeq = regexp.MustCompile("/+")

func newSystem(uiUri string, uiPath string, cookieMaxAge int, idpCont idpContainer) *system {
	uiUri = strings.TrimRight(uiUri, "/")
	uiUri = slashSeq.ReplaceAllString(uiUri, "/")
	if uiUri == "" {
		uiUri = "/html"
	}
	if uiUri[0] != '/' {
		uiUri = "/" + uiUri
	}
	log.Info("Use " + uiUri + " as UI uri")
	return &system{uiUri, uiPath, cookieMaxAge, idpCont}
}
