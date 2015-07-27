<!--
Copyright 2015 realglobe, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->


# TA 情報提供機能

UI 用に TA 情報を提供する。


## 1. 動作仕様

TA の指定は、TA の ID をパーセントエンコードし、パスにつなげて行う。

TA 情報は以下を最上位要素として含む JSON で返される。

* **`client_name`**
    * 名前。
      言語タグが付くことがある。


### 1.1. リクエスト例

```http
GET /api/info/ta/https%3A%2F%2Fta.example.org
Host: idp.example.org
```


### 1.2. レスポンス例

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
    "client_name#en": "That TA",
    "client_name#ja": "あの TA"
}
```
