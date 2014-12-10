package main

import (
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func init() {
	util.SetupConsoleLog("github.com/realglobe-Inc", level.OFF)
}

func TestBoot(t *testing.T) {
	// ////////////////////////////////
	// util.SetupConsoleLog("github.com/realglobe-Inc", level.ALL)
	// defer util.SetupConsoleLog("github.com/realglobe-Inc", level.OFF)
	// ////////////////////////////////

	port, err := util.FreePort()
	if err != nil {
		t.Fatal(err)
	}

	sys := &system{
		IdpLister: NewMemoryIdpLister(0),
	}
	go serve(sys, "tcp", "", port, "http")

	// サーバ起動待ち。
	time.Sleep(50 * time.Millisecond)

	req, err := http.NewRequest("GET", "http://localhost:"+strconv.Itoa(port)+listPagePath, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusAccepted {
		t.Error(resp)
	}
}
