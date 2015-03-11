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
	"encoding/json"
	logutil "github.com/realglobe-Inc/edo-lib/log"
	"github.com/realglobe-Inc/edo-lib/server"
	"github.com/realglobe-Inc/edo-lib/test"
	"github.com/realglobe-Inc/go-lib/erro"
	"github.com/realglobe-Inc/go-lib/rglog/level"
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

func setupIdpSelector() (sys *system, urlHead string, shutCh chan struct{}, err error) {
	port, err := test.FreePort()
	if err != nil {
		return nil, "", nil, erro.Wrap(err)
	}

	if sys == nil {
		sys = newTestSystem()
	}
	shutCh = make(chan struct{}, 10)
	urlHead = "http://localhost:" + strconv.Itoa(port)

	go serve(sys, "tcp", "", port, "http", shutCh)
	// 起動待ち。
	for i := time.Nanosecond; i < time.Second; i *= 2 {
		req, err := http.NewRequest("GET", urlHead+okPath, nil)
		if err != nil {
			os.RemoveAll(sys.uiPath)
			sys.close()
			shutCh <- struct{}{}
			return nil, "", nil, erro.Wrap(err)
		}
		req.Header.Set("Connection", "close")
		resp, err := (&http.Client{}).Do(req)
		if err != nil {
			// ちょっと待って再挑戦。
			time.Sleep(i)
			continue
		}
		// ちゃんとつながったので終わり。
		resp.Body.Close()
		return sys, urlHead, shutCh, nil
	}
	// 時間切れ。
	os.RemoveAll(sys.uiPath)
	sys.close()
	shutCh <- struct{}{}
	return nil, "", nil, erro.New("time out")
}

func TestSelectPage(t *testing.T) {
	// ////////////////////////////////
	// logutil.SetupConsole("github.com/realglobe-Inc", level.ALL)
	// defer logutil.SetupConsole("github.com/realglobe-Inc", level.OFF)
	// ////////////////////////////////

	sys, urlHead, shutCh, err := setupIdpSelector()
	if err != nil {
		t.Fatal(err)
	}
	defer sys.close()
	defer os.RemoveAll(sys.uiPath)
	defer func() { shutCh <- struct{}{} }()

	body := "<html><head><title>さんぷる</title></head><body>いろはに</body></html>"
	if err := ioutil.WriteFile(filepath.Join(sys.uiPath, "index.html"), []byte(body), filePerm); err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", urlHead+selectUri, nil)
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

	sys, urlHead, shutCh, err := setupIdpSelector()
	if err != nil {
		t.Fatal(err)
	}
	defer sys.close()
	defer os.RemoveAll(sys.uiPath)
	defer func() { shutCh <- struct{}{} }()

	idp := &idProvider{"https://example.com", "さんぷる", "https://example.com/login"}
	sys.idpCont.(*memoryIdpContainer).add(idp)

	req, err := http.NewRequest("GET", urlHead+listUri, nil)
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

	sys, urlHead, shutCh, err := setupIdpSelector()
	if err != nil {
		t.Fatal(err)
	}
	defer sys.close()
	defer os.RemoveAll(sys.uiPath)
	defer func() { shutCh <- struct{}{} }()

	idp := &idProvider{"https://example.com", "さんぷる", urlHead + sys.uiUri}
	sys.idpCont.(*memoryIdpContainer).add(idp)
	body := "<html><head><title>さんぷる</title></head><body>いろはに</body></html>"
	if err := ioutil.WriteFile(filepath.Join(sys.uiPath, "index.html"), []byte(body), filePerm); err != nil {
		t.Fatal(err)
	}

	// サーバ起動待ち。
	time.Sleep(50 * time.Millisecond)

	req, err := http.NewRequest("GET", urlHead+redirectUri+"?idp="+url.QueryEscape("https://example.com"), nil)
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
