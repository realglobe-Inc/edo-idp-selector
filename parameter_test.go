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
	"io/ioutil"
	"os"
	"testing"

	"github.com/realglobe-Inc/go-lib/rglog/level"
)

func TestParameters(t *testing.T) {
	file, err := ioutil.TempFile("", "edo-idp-selector")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	if _, err := file.Write([]byte(`
-consLv=ALL
-logType=logType
-logLv=ERR
-logPath=logPath
-logSize=16799
-logNum=63
-logAddr=logAddr
-logTag=logTag
-socType=socType
-socPath=socPath
-socPort=2440
-protType=protType
-pathOk=pathOk
-pathStart=pathStart
-pathSel=pathSel
-pathIdp=pathIdp
-pathUi=pathUi
-pathSelUi=pathSelUi
-uiDir=uiDir
-tmplErr=tmpErr
-sessLabel=sessLabel
-sessLen=69
-sessExpIn=24069s
-sessRefDelay=29212s
-sessDbExpIn=5806s
-ticLen=30
-ticExpIn=14770s
-redTimeout=21587s
-redPoolSize=68
-redPoolExpIn=5115s
-monTimeout=13469s
-webDbType=webDbType
-webDbAddr=webDbAddr
-webDbTag=webDbTag
-webDbExpIn=6949s
-idpDbType=idpDbType
-idpDbAddr=idpDbAddr
-idpDbTag=idpDbTag
-idpDbTag2=idpDbTag2
-taDbType=taDbType
-taDbAddr=taDbAddr
-taDbTag=taDbTag
-taDbTag2=taDbTag2
-sessDbType=sessDbType
-sessDbAddr=sessDbAddr
-sessDbTag=sessDbTag
-cookPath=cookPath
-cookSec=false
`)); err != nil {
		t.Fatal(err)
	} else if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	param, err := parseParameters("", "-c", file.Name())
	if err != nil {
		t.Fatal(err)
	} else if param.consLv != level.ALL {
		t.Error(param.consLv)
		t.Fatal(level.ALL)
	}
	// TODO 続きは後で。
}
