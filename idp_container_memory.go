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
	"github.com/realglobe-Inc/edo-lib/driver"
	"time"
)

type memoryIdpContainer idpContainerImpl

// スレッドセーフ。
func newMemoryIdpContainer(staleDur, expiDur time.Duration) *memoryIdpContainer {
	return (*memoryIdpContainer)(&idpContainerImpl{driver.NewMemoryListedKeyValueStore(staleDur, expiDur)})
}

func (this *memoryIdpContainer) get(idpId string) (*idProvider, error) {
	return ((*idpContainerImpl)(this)).get(idpId)
}

func (this *memoryIdpContainer) list(filter map[string]string) ([]*idProvider, error) {
	return ((*idpContainerImpl)(this)).list(filter)
}

func (this *memoryIdpContainer) close() error {
	return ((*idpContainerImpl)(this)).base.(driver.KeyValueStore).Close()
}

func (this *memoryIdpContainer) add(idp *idProvider) {
	((*idpContainerImpl)(this)).base.(driver.KeyValueStore).Put(idp.Id, idp)
}
