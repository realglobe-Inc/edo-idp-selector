package main

import (
	"github.com/realglobe-Inc/edo/driver"
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
