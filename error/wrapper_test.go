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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/realglobe-Inc/edo-lib/server"
	"github.com/realglobe-Inc/go-lib/erro"
)

func TestWrapPagePanic(t *testing.T) {
	hndl := WrapPage(server.NewStopper(), func(w http.ResponseWriter, r *http.Request) error {
		panic("abcde")
	}, nil)

	w := httptest.NewRecorder()
	hndl.ServeHTTP(w, &http.Request{})

	if w.Code != http.StatusInternalServerError {
		t.Error(w.Code)
		t.Fatal(http.StatusInternalServerError)
	} else if contType, contType2 := "text/html", w.HeaderMap.Get("Content-Type"); contType2 != contType {
		t.Error(contType2)
		t.Fatal(contType)
	}
}

func TestWrapPageError(t *testing.T) {
	hndl := WrapPage(server.NewStopper(), func(w http.ResponseWriter, r *http.Request) error {
		return erro.New("abcde")
	}, nil)

	w := httptest.NewRecorder()
	hndl.ServeHTTP(w, &http.Request{})

	if w.Code != http.StatusInternalServerError {
		t.Error(w.Code)
		t.Fatal(http.StatusInternalServerError)
	} else if contType, contType2 := "text/html", w.HeaderMap.Get("Content-Type"); contType2 != contType {
		t.Error(contType2)
		t.Fatal(contType)
	}
}

func TestWrapPageServerError(t *testing.T) {
	hndl := WrapPage(server.NewStopper(), func(w http.ResponseWriter, r *http.Request) error {
		return New("access_denied", "you are alone", http.StatusForbidden, nil)
	}, nil)

	w := httptest.NewRecorder()
	hndl.ServeHTTP(w, &http.Request{})

	if w.Code != http.StatusForbidden {
		t.Error(w.Code)
		t.Fatal(http.StatusForbidden)
	} else if contType, contType2 := "text/html", w.HeaderMap.Get("Content-Type"); contType2 != contType {
		t.Error(contType2)
		t.Fatal(contType)
	}
}

func TestWrapApiPanic(t *testing.T) {
	hndl := WrapApi(server.NewStopper(), func(w http.ResponseWriter, r *http.Request) error {
		panic("abcde")
	})

	w := httptest.NewRecorder()
	hndl.ServeHTTP(w, &http.Request{})

	if w.Code != http.StatusInternalServerError {
		t.Error(w.Code)
		t.Fatal(http.StatusInternalServerError)
	} else if contType, contType2 := "application/json", w.HeaderMap.Get("Content-Type"); contType2 != contType {
		t.Error(contType2)
		t.Fatal(contType)
	}
}

func TestWrapApiError(t *testing.T) {
	hndl := WrapApi(server.NewStopper(), func(w http.ResponseWriter, r *http.Request) error {
		return erro.New("abcde")
	})

	w := httptest.NewRecorder()
	hndl.ServeHTTP(w, &http.Request{})

	if w.Code != http.StatusInternalServerError {
		t.Error(w.Code)
		t.Fatal(http.StatusInternalServerError)
	} else if contType, contType2 := "application/json", w.HeaderMap.Get("Content-Type"); contType2 != contType {
		t.Error(contType2)
		t.Fatal(contType)
	}
}

func TestWrapApiServerError(t *testing.T) {
	hndl := WrapApi(server.NewStopper(), func(w http.ResponseWriter, r *http.Request) error {
		return New("access_denied", "you are alone", http.StatusForbidden, nil)
	})

	w := httptest.NewRecorder()
	hndl.ServeHTTP(w, &http.Request{})

	if w.Code != http.StatusForbidden {
		t.Error(w.Code)
		t.Fatal(http.StatusForbidden)
	} else if contType, contType2 := "application/json", w.HeaderMap.Get("Content-Type"); contType2 != contType {
		t.Error(contType2)
		t.Fatal(contType)
	}
}
