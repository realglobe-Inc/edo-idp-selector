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

package error

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib/erro"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestRedirectError(t *testing.T) {
	origErr := New("invalid_request", "invalid request", http.StatusBadRequest, nil)

	r, err := http.NewRequest("GET", "https://idp.example.org/", nil)
	if err != nil {
		t.Fatal(err)
	}
	uri, err := url.Parse("https://ta.example.org/callback")
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	RedirectError(w, r, erro.Wrap(origErr), uri, nil)

	if w.Code != http.StatusFound {
		t.Error(w.Code)
		t.Fatal(http.StatusFound)
	} else if uri2, err := url.Parse(w.HeaderMap.Get("Location")); err != nil {
		t.Fatal(err)
	} else if q := uri2.Query(); q.Get("error") != origErr.ErrorCode() {
		t.Error(q.Get("error"))
		t.Fatal(origErr.ErrorCode())
	} else if q.Get("error_description") != origErr.ErrorDescription() {
		t.Error(q.Get("error_description"))
		t.Fatal(origErr.ErrorDescription())
	}
}

func TestRespondJson(t *testing.T) {
	origErr := New("invalid_request", "invalid request", http.StatusBadRequest, nil)

	w := httptest.NewRecorder()
	RespondJson(w, nil, erro.Wrap(origErr), nil)

	if w.Code != origErr.Status() {
		t.Error(w.Code)
		t.Fatal(origErr.Status())
	} else if w.HeaderMap.Get("Content-Type") != "application/json" {
		t.Error(w.HeaderMap.Get("Content-Type"))
		t.Fatal("application/json")
	} else if w.Body == nil {
		t.Fatal("no body")
	}

	data, _ := ioutil.ReadAll(w.Body)
	var buff struct {
		Error             string
		Error_description string
	}
	if err := json.Unmarshal(data, &buff); err != nil {
		t.Fatal(err)
	} else if buff.Error != origErr.ErrorCode() {
		t.Error(buff.Error)
		t.Fatal(origErr.ErrorCode())
	} else if buff.Error_description != origErr.ErrorDescription() {
		t.Error(buff.Error_description)
		t.Fatal(origErr.ErrorDescription())
	}
}

func TestRespondHtml(t *testing.T) {
	origErr := New("invalid_request", "invalid request", http.StatusBadRequest, nil)

	w := httptest.NewRecorder()
	RespondHtml(w, nil, erro.Wrap(origErr), nil, nil)

	if w.Code != origErr.Status() {
		t.Error(w.Code)
		t.Fatal(origErr.Status())
	} else if w.HeaderMap.Get("Content-Type") != "text/html" {
		t.Error(w.HeaderMap.Get("Content-Type"))
		t.Fatal("text/html")
	} else if w.Body == nil {
		t.Fatal("no body")
	}
}

func TestRespondHtmlTemplate(t *testing.T) {
	origErr := New("invalid_request", "invalid request", http.StatusBadRequest, nil)

	file, err := ioutil.TempFile("", "edo-idp-selector.error")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	if _, err := file.Write([]byte("{{.Status}}")); err != nil {
		t.Fatal(err)
	}
	file.Close()

	tmpl, err := template.ParseFiles(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	RespondHtml(w, nil, erro.Wrap(origErr), tmpl, nil)

	if w.Code != origErr.Status() {
		t.Error(w.Code)
		t.Fatal(origErr.Status())
	} else if w.HeaderMap.Get("Content-Type") != "text/html" {
		t.Error(w.HeaderMap.Get("Content-Type"))
		t.Fatal("text/html")
	} else if w.Body == nil {
		t.Fatal("no body")
	}

	buff, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	} else if string(buff) != strconv.Itoa(origErr.Status()) {
		t.Error(string(buff))
		t.Fatal(origErr.Status())
	}
}

func TestRespondHtmlTemplateFunction(t *testing.T) {
	origErr := New("invalid_request", "invalid request", http.StatusBadRequest, nil)

	file, err := ioutil.TempFile("", "edo-lib")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	if _, err := file.Write([]byte("{{.Status}} {{.StatusText}} {{.Error}} {{.Description}} {{.Debug}}")); err != nil {
		t.Fatal(err)
	}
	file.Close()

	tmpl, err := template.ParseFiles(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	RespondHtml(w, nil, erro.Wrap(origErr), tmpl, nil)

	if w.Code != origErr.Status() {
		t.Error(w.Code)
		t.Fatal(origErr.Status())
	} else if w.HeaderMap.Get("Content-Type") != "text/html" {
		t.Error(w.HeaderMap.Get("Content-Type"))
		t.Fatal("text/html")
	} else if w.Body == nil {
		t.Fatal("no body")
	}

	buff, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	} else if parts := strings.Fields(string(buff)); parts[0] != strconv.Itoa(origErr.Status()) {
		t.Error(string(buff))
		t.Fatal(origErr.Status())
	}
}
