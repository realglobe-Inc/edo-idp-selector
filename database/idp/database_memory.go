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

package idp

import (
	"regexp"
	"strings"
)

// メモリ上の ID プロバイダ情報の格納庫。
type memoryDb struct {
	idToElem map[string]Element
}

func NewMemoryDb(elems []Element) Db {
	idToElem := map[string]Element{}
	for _, elem := range elems {
		idToElem[elem.Id()] = elem
	}
	return &memoryDb{
		idToElem,
	}
}

// 取得。
func (this *memoryDb) Get(id string) (Element, error) {
	return this.idToElem[id], nil
}

// まとめて取得。
func (this *memoryDb) Search(filter map[string]string) ([]Element, error) {
	regExps := map[string]*regexp.Regexp{}
	for tag, v := range filter {
		regExp, err := regexp.Compile(v)
		if err != nil {
			// こっちは悪くないのでエラーは返さない。
			return nil, nil
		}
		regExps[tag] = regExp
	}

	elems := []Element{}
	for _, elem := range this.idToElem {
		ok := true
		for tag, regExp := range regExps {
			switch tag {
			case tagIssuer:
				if !regExp.MatchString(elem.Id()) {
					ok = false
					break
				}
			case tagIssuer_name:
				if !regExp.MatchString(elem.Names()[""]) {
					ok = false
					break
				}
			case tagAuthorization_endpoint:
				if !regExp.MatchString(elem.AuthUri()) {
					ok = false
					break
				}
			case tagToken_endpoint:
				if !regExp.MatchString(elem.TokenUri()) {
					ok = false
					break
				}
			case tagUserinfo_endpoint:
				if !regExp.MatchString(elem.AccountUri()) {
					ok = false
					break
				}
			case tagCooperation_from_endpoint:
				if !regExp.MatchString(elem.CoopFromUri()) {
					ok = false
					break
				}
			case tagCooperation_to_endpoint:
				if !regExp.MatchString(elem.CoopToUri()) {
					ok = false
					break
				}
			default:
				if !strings.HasPrefix(tag, tagIssuer_name) {
					continue
				}

				// issuer_nameXXX
				// switch で調べてあるので issuer_name ではない。

				if tag[len(tagIssuer_name)] != '#' {
					continue
				}

				if !regExp.MatchString(elem.Names()[tag[len(tagIssuer_name)+1:]]) {
					ok = false
					break
				}
			}
		}
		if ok {
			elems = append(elems, elem)
		}
	}
	return elems, nil
}
