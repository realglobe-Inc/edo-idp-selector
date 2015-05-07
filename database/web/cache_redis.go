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

import (
	"github.com/garyburd/redigo/redis"
	"github.com/realglobe-Inc/go-lib/erro"
	"time"
)

// redis による自身の鍵のキャッシュ。
type redisCache struct {
	base  Db
	pool  *redis.Pool
	tag   string
	expIn time.Duration
}

func NewRedisCache(base Db, pool *redis.Pool, tag string, expIn time.Duration) Db {
	return &redisCache{
		base:  base,
		pool:  pool,
		tag:   tag,
		expIn: expIn,
	}
}

func (this *redisCache) Get(uri string) (Element, error) {
	conn := this.pool.Get()
	defer conn.Close()

	if data, err := redis.Bytes(conn.Do("GET", this.tag+uri)); err != nil {
		if err != redis.ErrNil {
			log.Warn(erro.Wrap(err))
			// キャッシュが取れなくても諦めない。
		}
	} else {
		// キャッシュされてた。
		return newElement(uri, data), nil
	}

	// キャッシュされてなかった。
	elem, err := this.base.Get(uri)
	if err != nil {
		return nil, erro.Wrap(err)
	} else if elem == nil {
		return nil, nil
	}

	// キャッシュする。
	if _, err := conn.Do("SET", this.tag+uri, elem.Data(), "PX", int64(this.expIn/time.Millisecond)); err != nil {
		log.Warn(erro.Wrap(err))
		// キャッシュできなくても諦めない。
	}

	return elem, nil
}
