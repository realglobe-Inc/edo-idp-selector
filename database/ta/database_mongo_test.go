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

package ta

import (
	"strconv"
	"testing"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// テストするなら、mongodb を立てる必要あり。
// 立ってなかったらテストはスキップ。
var monPool, _ = mgo.DialWithTimeout("localhost", time.Minute)

func init() {
	if monPool != nil {
		monPool.SetSyncTimeout(time.Minute)
	}
}

const (
	test_coll = "test-collection"
)

func TestMongoDb(t *testing.T) {
	if monPool == nil {
		t.SkipNow()
	}

	test_db := "test-db-" + strconv.FormatInt(time.Now().UnixNano(), 16)
	conn := monPool.New()
	defer conn.Close()
	if err := conn.DB(test_db).C(test_coll).Insert(bson.M{
		"client_id":      test_elem.Id(),
		"client_name":    test_elem.Names()[""],
		"client_name#ja": test_elem.Names()["ja"],
		"redirect_uris":  []string{test_rediUri},
		"jwks": []bson.M{
			{
				"kty": "EC",
				"crv": "P-256",
				"x":   "lpHYO1qpjU95B2sThPR2-1jv44axgaEDkQtcKNE-oZs",
				"y":   "soy5O11SFFFeYdhQVodXlYPIpeo0pCS69IxiVPPf0Tk",
				"d":   "3BhkCluOkm8d8gvaPD5FDG2zeEw2JKf3D5LwN-mYmsw",
			},
		},
		"pairwise": func() string {
			if test_elem.Pairwise() {
				return "pairwise"
			} else {
				return "public"
			}
		}(),
		"sector_identifier_uri": test_sect,
	}); err != nil {
		t.Fatal(err)
	}
	defer conn.DB(test_db).DropDatabase()

	testDb(t, NewMongoDb(monPool, test_db, test_coll, nil))
}
