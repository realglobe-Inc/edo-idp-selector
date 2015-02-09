package main

import (
	"flag"
	"fmt"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type parameters struct {
	// 画面ログ表示重要度。
	consLv level.Level

	// 追加ログ種別。
	logType string
	// 追加ログ表示重要度。
	logLv level.Level
	// ログファイルパス。
	logPath string
	// fluentd アドレス。
	fluAddr string
	// fluentd 用タグ。
	fluTag string

	// ソケット種別。
	socType string
	// UNIX ソケット。
	socPath string
	// TCP ソケット。
	socPort int

	// プロトコル種別。
	protType string

	// キャッシュを最新とみなす期間。
	caStaleDur time.Duration
	// キャッシュを廃棄するまでの期間。
	caExpiDur time.Duration

	// UI 用 HTML を提供する URI。
	uiUri string
	// UI 用 HTML を置くディレクトリパス。
	uiPath string

	// IdP 格納庫種別。
	idpContType string
	// IdP 格納庫ディレクトリパス。
	idpContPath string
	// IdP 格納庫 mongodb アドレス。
	idpContUrl string
	// IdP 格納庫 mongodb データベース名。
	idpContDb string
	// IdP 格納庫 mongodb コレクション名。
	idpContColl string

	// cookie の有効期間（秒）。
	cookMaxAge int
}

func parseParameters(args ...string) (param *parameters, err error) {

	const label = "edo-idp-selector"

	flags := flag.NewFlagSet(label+" parameters", flag.ExitOnError)
	flags.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, "  "+args[0]+" [{FLAG}...]")
		fmt.Fprintln(os.Stderr, "FLAG:")
		flags.PrintDefaults()
	}

	param = &parameters{}

	flags.Var(level.Var(&param.consLv, level.INFO), "consLv", "Console log level.")
	flags.StringVar(&param.logType, "logType", "", "Extra log type. file/fluentd")
	flags.Var(level.Var(&param.logLv, level.ALL), "logLv", "Extra log level.")
	flags.StringVar(&param.logPath, "logPath", filepath.Join(filepath.Dir(os.Args[0]), "log", label+".log"), "File log path.")
	flags.StringVar(&param.fluAddr, "fluAddr", "localhost:24224", "fluentd address.")
	flags.StringVar(&param.fluTag, "fluTag", "edo."+label, "fluentd tag.")

	flags.StringVar(&param.socType, "socType", "tcp", "Socket type. tcp/unix")
	flags.StringVar(&param.socPath, "socPath", filepath.Join(filepath.Dir(os.Args[0]), "run", label+".soc"), "UNIX socket path.")
	flags.IntVar(&param.socPort, "socPort", 16030, "TCP socket port.")

	flags.StringVar(&param.protType, "protType", "http", "Protocol type. http/fcgi")

	flags.DurationVar(&param.caStaleDur, "caStaleDur", 5*time.Minute, "Cache fresh duration.")
	flags.DurationVar(&param.caExpiDur, "caExpiDur", 30*time.Minute, "Cache expiration duration.")

	flags.StringVar(&param.uiUri, "uiUri", "/html", "UI uri.")
	flags.StringVar(&param.uiPath, "uiPath", filepath.Join(filepath.Dir(os.Args[0]), "html"), "Protocol type. http/fcgi")

	flags.StringVar(&param.idpContType, "idpContType", "file", "IdP container type.")
	flags.StringVar(&param.idpContPath, "idpContPath", filepath.Join(filepath.Dir(os.Args[0]), "idps"), "IdP container directory.")
	flags.StringVar(&param.idpContUrl, "idpContUrl", "localhost", "IdP container address.")
	flags.StringVar(&param.idpContDb, "idpContDb", "edo", "IdP container database name.")
	flags.StringVar(&param.idpContColl, "idpContColl", "ta_uris", "IdP container collection name.")

	flags.IntVar(&param.cookMaxAge, "cookMaxAge", 7*24*60*60, "Cookie expiration duration (second).")

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
			flags.Parse(strings.Fields(string(buff)))
		}
	}
	flags.Parse(args[1:])

	if l := len(flags.Args()); l > 0 {
		log.Warn("Ignore extra parameters ", flags.Args(), ".")
	}

	return param, nil
}
