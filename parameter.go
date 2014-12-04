package main

import (
	"flag"
	"fmt"
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type parameters struct {
	// 画面表示ログ。
	consLv level.Level

	// 追加ログ。
	logType string
	logLv   level.Level

	// ファイルログ。
	logPath string

	// fluentd ログ。
	fluAddr string
	fluTag  string

	// ID プロバイダリスト。
	idpListType string

	// ファイルベース ID プロバイダリスト。
	idpListPath string

	// Web ベース ID プロバイダリスト。
	idpListAddr string

	// mongo ID プロバイダリスト。
	idpListUrl  string
	idpListDb   string
	idpListColl string

	// ID プロバイダ属性レジストリ。
	idpAttrRegType string

	// ファイルベース ID プロバイダ属性レジストリ。
	idpAttrRegPath string

	// Web ベース ID プロバイダ属性レジストリ。
	idpAttrRegAddr string

	// mongo ID プロバイダ属性レジストリ。
	idpAttrRegUrl  string
	idpAttrRegDb   string
	idpAttrRegColl string

	// ソケット。
	socType string

	// UNIX ソケット。
	socPath string

	// TCP ソケット。
	socPort int

	// プロトコル。
	protType string

	// cookie の有効期間（秒）。
	cookieMaxAge int
}

func parseParameters(args ...string) (param *parameters, err error) {

	const label = "edo-idp-selector"

	flags := util.NewFlagSet(label+" parameters", flag.ExitOnError)
	flags.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, "  "+args[0]+" [{FLAG}...]")
		fmt.Fprintln(os.Stderr, "FLAG:")
		flags.PrintDefaults()
	}

	param = &parameters{}

	flags.Var(level.Var(&param.consLv, level.INFO), "consLv", "Console log level.")
	flags.StringVar(&param.logType, "logType", "", "Extra log type.")
	flags.Var(level.Var(&param.logLv, level.ALL), "logLv", "Extra log level.")
	flags.StringVar(&param.logPath, "logPath", filepath.Join(filepath.Dir(os.Args[0]), "log", label+".log"), "File log path.")
	flags.StringVar(&param.fluAddr, "fluAddr", "localhost:24224", "fluentd address.")
	flags.StringVar(&param.fluTag, "fluTag", "edo."+label, "fluentd tag.")

	flags.StringVar(&param.idpListType, "idpListType", "web", "ID provider lister type.")
	flags.StringVar(&param.idpListPath, "idpListPath", filepath.Join(filepath.Dir(os.Args[0]), "idps"), "ID provider lister directory.")
	flags.StringVar(&param.idpListAddr, "idpListAddr", "http://localhost:16031", "ID provider lister address.")
	flags.StringVar(&param.idpListUrl, "idpListUrl", "localhost", "ID provider lister address.")
	flags.StringVar(&param.idpListDb, "idpListDb", "edo", "ID provider lister database name.")
	flags.StringVar(&param.idpListColl, "idpListColl", "idps", "ID provider lister collection name.")

	flags.StringVar(&param.idpAttrRegType, "idpAttrRegType", "web", "ID provider attribute provider type.")
	flags.StringVar(&param.idpAttrRegPath, "idpAttrRegPath", filepath.Join(filepath.Dir(os.Args[0]), "idp_attributes"), "ID provider attribute provider directory.")
	flags.StringVar(&param.idpAttrRegAddr, "idpAttrRegAddr", "http://localhost:16032", "ID provider attribute provider address.")
	flags.StringVar(&param.idpAttrRegUrl, "idpAttrRegUrl", "localhost", "ID provider attribute provider address.")
	flags.StringVar(&param.idpAttrRegDb, "idpAttrRegDb", "edo", "ID provider attribute provider database name.")
	flags.StringVar(&param.idpAttrRegColl, "idpAttrRegColl", "idp_attributes", "ID provider attribute provider collection name.")

	flags.StringVar(&param.socType, "socType", "tcp", "Socket type.")
	flags.StringVar(&param.socPath, "socPath", filepath.Join(filepath.Dir(os.Args[0]), "run", label+".soc"), "UNIX socket path.")
	flags.IntVar(&param.socPort, "socPort", 16030, "TCP socket port.")

	flags.StringVar(&param.protType, "protType", "http", "Protocol type.")

	flags.IntVar(&param.cookieMaxAge, "cookieMaxAge", 7*24*60*60, "Cookie expiration duration (second).")

	var config string
	flags.StringVar(&config, "f", "", "Config file path.")

	// 実行引数を読んで、設定ファイルを指定させてから、
	// 設定ファイルを読んで、また実行引数を読む。
	flags.Parse(args[1:])
	if config != "" {
		if buff, err := ioutil.ReadFile(config); err != nil {
			if !os.IsNotExist(err) {
				return nil, erro.Wrap(err)
			}
			log.Warn("Config file " + config + " is not exist.")
		} else {
			flags.CompleteParse(strings.Fields(string(buff)))
		}
	}
	flags.Parse(args[1:])

	if l := len(flags.Args()); l > 0 {
		log.Warn("Ignore extra parameters ", flags.Args(), ".")
	}

	return param, nil
}
