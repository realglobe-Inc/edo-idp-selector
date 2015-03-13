// Copyright 2015 realglobe, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	logutil "github.com/realglobe-Inc/edo-lib/log"
	"github.com/realglobe-Inc/edo-lib/server"
	"github.com/realglobe-Inc/go-lib/erro"
	"github.com/realglobe-Inc/go-lib/rglog"
	"net/http"
	"os"
)

func main() {
	var exitCode = 0
	defer func() {
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	}()
	defer rglog.Flush()

	logutil.InitConsole("github.com/realglobe-Inc")

	param, err := parseParameters(os.Args...)
	if err != nil {
		log.Err(erro.Unwrap(err))
		log.Debug(err)
		exitCode = 1
		return
	}

	logutil.SetupConsole("github.com/realglobe-Inc", param.consLv)
	if err := logutil.Setup("github.com/realglobe-Inc", param.logType, param.logLv, param); err != nil {
		log.Err(erro.Unwrap(err))
		log.Debug(err)
		exitCode = 1
		return
	}

	if err := mainCore(param); err != nil {
		err = erro.Wrap(err)
		log.Err(erro.Unwrap(err))
		log.Debug(err)
		exitCode = 1
		return
	}

	log.Info("Shut down.")
}

// system を準備する。
func mainCore(param *parameters) error {
	var idpCont idpContainer
	switch param.idpContType {
	case "file":
		idpCont = newFileIdpContainer(param.idpContPath, param.caStaleDur, param.caExpiDur)
		log.Info("Use file IdP container in " + param.idpContPath)
	default:
		return erro.New("invalid IdP container type " + param.idpContType)
	}

	sys := newSystem(
		param.uiUri,
		param.uiPath,
		param.cookMaxAge,
		idpCont,
	)
	defer sys.close()
	return serve(sys, param.socType, param.socPath, param.socPort, param.protType, nil)
}

// 振り分ける。
const (
	selectUri   = "/"
	listUri     = "/list"
	redirectUri = "/redirect"
	okPath      = "/ok"
)

func serve(sys *system, socType, socPath string, socPort int, protType string, shutCh chan struct{}) error {
	routes := map[string]server.HandlerFunc{
		selectUri: func(w http.ResponseWriter, r *http.Request) error {
			return selectPage(sys, w, r)
		},
		listUri: func(w http.ResponseWriter, r *http.Request) error {
			return listApi(sys, w, r)
		},
		redirectUri: func(w http.ResponseWriter, r *http.Request) error {
			return redirectPage(sys, w, r)
		},
		okPath: func(w http.ResponseWriter, r *http.Request) error {
			return nil
		},
	}
	fileHndl := http.StripPrefix(sys.uiUri, http.FileServer(http.Dir(sys.uiPath)))
	for _, uri := range []string{sys.uiUri, sys.uiUri + "/"} {
		routes[uri] = func(w http.ResponseWriter, r *http.Request) error {
			fileHndl.ServeHTTP(w, r)
			return nil
		}
	}
	return server.TerminableServe(socType, socPath, socPort, protType, routes, shutCh, server.PanicErrorWrapper)
}
