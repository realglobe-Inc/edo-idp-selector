package main

import (
	logutil "github.com/realglobe-Inc/edo/util/log"
	"github.com/realglobe-Inc/edo/util/server"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/rglog"
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
	return serve(sys, param.socType, param.socPath, param.socPort, param.protType, nil)
}

// 振り分ける。
const (
	selectUri   = "/"
	listUri     = "/list"
	redirectUri = "/redirect"
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
