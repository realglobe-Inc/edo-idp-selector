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

package web

import ()

// web データの実装。
type element struct {
	uri  string
	data []byte
}

// 主にテスト用。
func New(uri string, data []byte) Element {
	return newElement(uri, data)
}

func newElement(uri string, data []byte) *element {
	return &element{uri, data}
}

func (this *element) Uri() string {
	return this.uri
}

func (this *element) Data() []byte {
	return this.data
}
