package main

import (
	"encoding/json"
	"github.com/realglobe-Inc/edo/util"
	logutil "github.com/realglobe-Inc/edo/util/log"
	"github.com/realglobe-Inc/edo/util/server"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func init() {
	logutil.SetupConsole("github.com/realglobe-Inc", level.OFF)
}

func newTestSystem() *system {
	uiPath, err := ioutil.TempDir("", "edo-idp-selector")
	if err != nil {
		panic(err)
	}
	return &system{
		uiUri:      "/html",
		uiPath:     uiPath,
		cookMaxAge: 1,
		idpCont:    newMemoryIdpContainer(0, 0),
	}
}

func TestSelectPage(t *testing.T) {
	// ////////////////////////////////
	// logutil.SetupConsole("github.com/realglobe-Inc", level.ALL)
	// defer logutil.SetupConsole("github.com/realglobe-Inc", level.OFF)
	// ////////////////////////////////

	port, err := util.FreePort()
	if err != nil {
		t.Fatal(err)
	}

	sys := newTestSystem()
	defer os.RemoveAll(sys.uiPath)
	shutCh := make(chan struct{}, 10)
	defer func() { shutCh <- struct{}{} }()
	go serve(sys, "tcp", "", port, "http", shutCh)
	body := "<html><head><title>さんぷる</title></head><body>いろはに</body></html>"
	if err := ioutil.WriteFile(filepath.Join(sys.uiPath, "index.html"), []byte(body), filePerm); err != nil {
		t.Fatal(err)
	}

	// サーバ起動待ち。
	time.Sleep(50 * time.Millisecond)

	req, err := http.NewRequest("GET", "http://localhost:"+strconv.Itoa(port)+selectUri, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Error(resp)
	} else if buff, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fatal(err)
	} else if string(buff) != body {
		t.Error(string(buff))
	}
}

func TestListPage(t *testing.T) {
	// ////////////////////////////////
	// logutil.SetupConsole("github.com/realglobe-Inc", level.ALL)
	// defer logutil.SetupConsole("github.com/realglobe-Inc", level.OFF)
	// ////////////////////////////////

	port, err := util.FreePort()
	if err != nil {
		t.Fatal(err)
	}

	sys := newTestSystem()
	defer os.RemoveAll(sys.uiPath)
	shutCh := make(chan struct{}, 10)
	defer func() { shutCh <- struct{}{} }()
	go serve(sys, "tcp", "", port, "http", shutCh)

	idp := &idProvider{"https://example.com", "さんぷる", "https://example.com/login"}
	sys.idpCont.(*memoryIdpContainer).add(idp)

	// サーバ起動待ち。
	time.Sleep(50 * time.Millisecond)

	req, err := http.NewRequest("GET", "http://localhost:"+strconv.Itoa(port)+listUri, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Error(resp)
	}
	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	var idps []*idProvider
	if err := json.Unmarshal(buff, &idps); err != nil {
		t.Fatal(err)
	} else if len(idps) != 1 || !reflect.DeepEqual(idps[0], idp) {
		t.Error(idps)
	}
}

func TestRedirectPage(t *testing.T) {
	// ////////////////////////////////
	// logutil.SetupConsole("github.com/realglobe-Inc", level.ALL)
	// defer logutil.SetupConsole("github.com/realglobe-Inc", level.OFF)
	// ////////////////////////////////

	port, err := util.FreePort()
	if err != nil {
		t.Fatal(err)
	}

	sys := newTestSystem()
	defer os.RemoveAll(sys.uiPath)
	shutCh := make(chan struct{}, 10)
	defer func() { shutCh <- struct{}{} }()
	go serve(sys, "tcp", "", port, "http", shutCh)
	idp := &idProvider{"https://example.com", "さんぷる", "http://localhost:" + strconv.Itoa(port) + sys.uiUri}
	sys.idpCont.(*memoryIdpContainer).add(idp)
	body := "<html><head><title>さんぷる</title></head><body>いろはに</body></html>"
	if err := ioutil.WriteFile(filepath.Join(sys.uiPath, "index.html"), []byte(body), filePerm); err != nil {
		t.Fatal(err)
	}

	// サーバ起動待ち。
	time.Sleep(50 * time.Millisecond)

	req, err := http.NewRequest("GET", "http://localhost:"+strconv.Itoa(port)+redirectUri+"?idp="+url.QueryEscape("https://example.com"), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	server.LogResponse(level.ERR, resp, true)

	if resp.StatusCode != http.StatusOK {
		t.Error(resp)
	} else if buff, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fatal(err)
	} else if string(buff) != body {
		t.Error(string(buff))
	}
}
