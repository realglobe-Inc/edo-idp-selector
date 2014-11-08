package main

import (
	"github.com/realglobe-Inc/edo/driver"
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/rglog"
	"net/http"
	"os"
)

var exitCode = 0

func exit() {
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

func main() {
	defer exit()
	defer rglog.Flush()

	hndl := util.InitLog("github.com/realglobe-Inc")

	param, err := parseParameters(os.Args...)
	if err != nil {
		log.Err(erro.Unwrap(err))
		log.Debug(err)
		exitCode = 1
		return
	}

	hndl.SetLevel(param.consLv)
	if err := util.SetupLog("github.com/realglobe-Inc", param.logType, param.logLv, param.logPath, param.fluAddr, param.fluTag); err != nil {
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
	var err error

	var idpAttrReg driver.IdpAttributeProvider
	switch param.idpAttrRegType {
	case "file":
		idpAttrReg = driver.NewFileIdpAttributeProvider(param.idpAttrRegPath, 0)
		log.Info("Use file ID provider attribute provider " + param.idpAttrRegPath + ".")
	case "web":
		idpAttrReg = driver.NewWebIdpAttributeProvider(param.idpAttrRegAddr)
		log.Info("Use web ID provider attribute provider " + param.idpAttrRegAddr + ".")
	case "mongo":
		idpAttrReg, err = driver.NewMongoIdpAttributeProvider(param.idpAttrRegUrl, param.idpAttrRegDb, param.idpAttrRegColl, 0)
		if err != nil {
			return erro.Wrap(err)
		}
		log.Info("Use mongodb ID provider attribute provider " + param.idpAttrRegUrl + ".")
	default:
		return erro.New("invalid ID provider attribute provider type " + param.idpAttrRegType + ".")
	}

	var idpList driver.IdpLister
	switch param.idpListType {
	case "file":
		idpList = driver.NewFileIdpLister(param.idpListPath, 0)
		log.Info("Use file ID provider lister " + param.idpListPath + ".")
	case "web":
		idpList = driver.NewWebIdpLister(param.idpListAddr)
		log.Info("Use web ID provider lister " + param.idpListAddr + ".")
	case "mongo":
		idpList, err = driver.NewMongoIdpLister(param.idpListUrl, param.idpListDb, param.idpListColl, 0)
		if err != nil {
			return erro.Wrap(err)
		}
		log.Info("Use mongodb ID provider lister " + param.idpListUrl + ".")
	default:
		return erro.New("invalid ID provider lister type " + param.idpListType + ".")
	}

	sys := &system{
		idpList,
		idpAttrReg,
		param.cookieMaxAge,
	}
	return serve(sys, param.socType, param.socPath, param.socPort, param.protType)
}

// 振り分ける。
const (
	routePagePath     = "/"
	listPagePath      = "/list"
	setCookiePagePath = "/set_cookie"
)

func serve(sys *system, socType, socPath string, socPort int, protType string) error {
	routes := map[string]util.HandlerFunc{
		routePagePath: func(w http.ResponseWriter, r *http.Request) error {
			return routePage(sys, w, r)
		},
		listPagePath: func(w http.ResponseWriter, r *http.Request) error {
			return listPage(sys, w, r)
		},
		setCookiePagePath: func(w http.ResponseWriter, r *http.Request) error {
			return setCookiePage(sys, w, r)
		},
	}
	return util.Serve(socType, socPath, socPort, protType, routes)
}
