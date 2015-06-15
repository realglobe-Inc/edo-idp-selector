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
	idpapi "github.com/realglobe-Inc/edo-idp-selector/api/idp"
	idpdb "github.com/realglobe-Inc/edo-idp-selector/database/idp"
	"github.com/realglobe-Inc/edo-idp-selector/database/session"
	tadb "github.com/realglobe-Inc/edo-idp-selector/database/ta"
	webdb "github.com/realglobe-Inc/edo-idp-selector/database/web"
	idperr "github.com/realglobe-Inc/edo-idp-selector/error"
	"github.com/realglobe-Inc/edo-idp-selector/page/idpselect"
	"github.com/realglobe-Inc/edo-lib/driver"
	logutil "github.com/realglobe-Inc/edo-lib/log"
	"github.com/realglobe-Inc/edo-lib/rand"
	"github.com/realglobe-Inc/edo-lib/server"
	"github.com/realglobe-Inc/go-lib/erro"
	"github.com/realglobe-Inc/go-lib/rglog"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"
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
		log.Debug(erro.Wrap(err))
		exitCode = 1
		return
	}

	logutil.SetupConsole("github.com/realglobe-Inc", param.consLv)
	if err := logutil.Setup("github.com/realglobe-Inc", param.logType, param.logLv, param); err != nil {
		log.Err(erro.Unwrap(err))
		log.Debug(erro.Wrap(err))
		exitCode = 1
		return
	}

	if err := serve(param); err != nil {
		log.Err(erro.Unwrap(err))
		log.Debug(erro.Wrap(err))
		exitCode = 1
		return
	}

	log.Info("Shut down")
}

func serve(param *parameters) (err error) {
	// バックエンドの準備。

	redPools := driver.NewRedisPoolSet(param.redTimeout, param.redPoolSize, param.redPoolExpIn)
	defer redPools.Close()
	monPools := driver.NewMongoPoolSet(param.monTimeout)
	defer monPools.Close()

	// web データ。
	var webDb webdb.Db
	switch param.webDbType {
	case "direct":
		webDb = webdb.NewDirectDb()
		log.Info("Get web data directly")
	case "redis":
		webDb = webdb.NewRedisCache(webdb.NewDirectDb(), redPools.Get(param.webDbAddr), param.webDbTag, param.webDbExpIn)
		log.Info("Get web data with redis " + param.webDbAddr + "<" + param.webDbTag + ">")
	default:
		return erro.New("invalid web data DB type " + param.webDbType)
	}

	// IdP 情報。
	var idpDb idpdb.Db
	switch param.idpDbType {
	case "mongo":
		pool, err := monPools.Get(param.idpDbAddr)
		if err != nil {
			return erro.Wrap(err)
		}
		idpDb = idpdb.NewMongoDb(pool, param.idpDbTag, param.idpDbTag2, webDb)
		log.Info("Use IdP info in mongodb " + param.idpDbAddr + "<" + param.idpDbTag + "." + param.idpDbTag2 + ">")
	default:
		return erro.New("invalid IdP DB type " + param.idpDbType)
	}

	// TA 情報。
	var taDb tadb.Db
	switch param.taDbType {
	case "mongo":
		pool, err := monPools.Get(param.taDbAddr)
		if err != nil {
			return erro.Wrap(err)
		}
		taDb = tadb.NewMongoDb(pool, param.taDbTag, param.taDbTag2, webDb)
		log.Info("Use TA info in mongodb " + param.taDbAddr + "<" + param.taDbTag + "." + param.taDbTag2 + ">")
	default:
		return erro.New("invalid TA DB type " + param.taDbType)
	}

	// セッション。
	var sessDb session.Db
	switch param.sessDbType {
	case "memory":
		sessDb = session.NewMemoryDb()
		log.Info("Save sessions in memory")
	case "redis":
		sessDb = session.NewRedisDb(redPools.Get(param.sessDbAddr), param.sessDbTag)
		log.Info("Save sessions in redis " + param.sessDbAddr + "<" + param.sessDbTag + ">")
	default:
		return erro.New("invalid session DB type " + param.sessDbType)
	}

	var errTmpl *template.Template
	if param.tmplErr != "" {
		errTmpl, err = template.ParseFiles(param.tmplErr)
		if err != nil {
			return erro.Wrap(err)
		}
	}

	idGen := rand.New(time.Minute)

	// バックエンドの準備完了。

	s := server.NewStopper()
	defer func() {
		// 処理の終了待ち。
		s.Lock()
		defer s.Unlock()
		for s.Stopped() {
			s.Wait()
		}
	}()

	selPage := idpselect.New(
		s,
		param.pathSelUi,
		errTmpl,
		param.sessLabel,
		param.sessLen,
		param.sessExpIn,
		param.sessRefDelay,
		param.sessDbExpIn,
		param.ticLen,
		param.ticExpIn,
		idpDb,
		taDb,
		sessDb,
		idGen,
		param.cookPath,
		param.cookSec,
		param.debug,
	)

	mux := http.NewServeMux()
	routes := map[string]bool{}
	mux.HandleFunc(param.pathOk, idperr.WrapPage(s, func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}, errTmpl))
	routes[param.pathOk] = true
	mux.HandleFunc(param.pathStart, selPage.HandleStart)
	routes[param.pathStart] = true
	mux.HandleFunc(param.pathSel, selPage.HandleSelect)
	routes[param.pathSel] = true
	mux.Handle(param.pathIdp, idpapi.New(s, idpDb, param.debug))
	routes[param.pathIdp] = true
	if param.uiDir != "" {
		// ファイル配信も自前でやる。
		pathUi := strings.TrimRight(param.pathUi, "/") + "/"
		mux.Handle(pathUi, http.StripPrefix(pathUi, http.FileServer(http.Dir(param.uiDir))))
		routes[param.pathUi] = true
	}

	if !routes["/"] {
		mux.HandleFunc("/", idperr.WrapPage(s, func(w http.ResponseWriter, r *http.Request) error {
			return erro.Wrap(idperr.New(idperr.Invalid_request, "invalid endpoint", http.StatusNotFound, nil))
		}, errTmpl))
	}

	return server.Serve(param, mux)
}
