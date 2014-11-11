package main

import (
	"flag"
	"fmt"
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/file"
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
	dsLogPath string

	// fluentd ログ。
	fluAddr  string
	dsFluTag string

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
	dsSocType string

	// UNIX ソケット。
	dsSocPath string

	// TCP ソケット。
	dsSocPort int

	// プロトコル。
	dsProtType string

	// cookie の有効期間（秒）。
	cookieMaxAge int
}

func parseParameters(args ...string) (param *parameters, err error) {

	const label = "idp-selector"

	flags := util.NewFlagSet("edo-"+label+" parameters", flag.ExitOnError)
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
	flags.StringVar(&param.dsLogPath, "dsLogPath", filepath.Join(os.TempDir(), "edo-"+label+".log"), "File log path.")
	flags.StringVar(&param.fluAddr, "fluAddr", "localhost:24224", "fluentd address.")
	flags.StringVar(&param.dsFluTag, "dsFluTag", "edo."+label, "fluentd tag.")

	flags.StringVar(&param.idpListType, "idpListType", "web", "ID provider lister type.")
	flags.StringVar(&param.idpListPath, "idpListPath", filepath.Join("sandbox", "idp-lister"), "ID provider lister directory.")
	flags.StringVar(&param.idpListAddr, "idpListAddr", "http://localhost:16031", "ID provider lister address.")
	flags.StringVar(&param.idpListUrl, "idpListUrl", "localhost", "ID provider lister address.")
	flags.StringVar(&param.idpListDb, "idpListDb", "edo", "ID provider lister database name.")
	flags.StringVar(&param.idpListColl, "idpListColl", "idp-lister", "ID provider lister collection name.")

	flags.StringVar(&param.idpAttrRegType, "idpAttrRegType", "web", "ID provider attribute registry type.")
	flags.StringVar(&param.idpAttrRegPath, "idpAttrRegPath", filepath.Join("sandbox", "id-provider-attribute-registry"), "ID provider attribute registry directory.")
	flags.StringVar(&param.idpAttrRegAddr, "idpAttrRegAddr", "http://localhost:9001", "ID provider attribute registry address.")
	flags.StringVar(&param.idpAttrRegUrl, "idpAttrRegUrl", "localhost", "ID provider attribute registry address.")
	flags.StringVar(&param.idpAttrRegDb, "idpAttrRegDb", "edo", "ID provider attribute registry database name.")
	flags.StringVar(&param.idpAttrRegColl, "idpAttrRegColl", "id-provider-attribute-registry", "ID provider attribute registry collection name.")

	flags.StringVar(&param.dsSocType, "dsSocType", "tcp", "Socket type.")
	flags.StringVar(&param.dsSocPath, "dsSocPath", filepath.Join(os.TempDir(), "edo-"+label), "UNIX socket path.")
	flags.IntVar(&param.dsSocPort, "dsSocPort", 16030, "TCP socket port.")

	flags.StringVar(&param.dsProtType, "dsProtType", "http", "Protocol type.")

	flags.IntVar(&param.cookieMaxAge, "cookieMaxAge", 7*24*60*60, "Cookie expiration duration (second).")

	var config string
	flags.StringVar(&config, "f", "", "Config file path.")

	// 実行引数を読んで、設定ファイルを指定させてから、
	// 設定ファイルを読んで、また実行引数を読む。
	flags.Parse(args[1:])
	if config != "" {
		if exist, err := file.IsExist(config); err != nil {
			return nil, erro.Wrap(err)
		} else if !exist {
			log.Warn("Config file " + config + " is not exist.")
		} else {
			buff, err := ioutil.ReadFile(config)
			if err != nil {
				return nil, erro.Wrap(err)
			}
			flags.CompleteParse(strings.Fields(string(buff)))
		}
	}
	flags.Parse(args[1:])

	if l := len(flags.Args()); l > 0 {
		log.Warn("Ignore extra parameters ", flags.Args(), ".")
	}

	return param, nil
}
