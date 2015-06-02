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


# ID プロバイダ情報提供機能

UI 用に ID プロバイダ情報を提供する。


## 1. 動作仕様

### 1.1. ID プロバイダ列挙エンドポイント

ID プロバイダ情報を列挙する。
クエリで絞り込める。
クエリは以下の形式の連言として用いる。

```
<タグ名>=<該当する値の正規表現>
```

レスポンスは [OpenID Connect Discovery 1.0 Section 4.2] 形式の ID プロバイダ情報の JSON 配列である。
ただし、全ての情報が返されるわけではない。


#### 1.1.1. リクエスト例

```http
GET /api/info/issuer?issuer=%5C.example%5C.org%24
Host: selector.example.org
```


#### 1.1.2. レスポンス例

```http
HTTP/1.1 200 OK
Content-Type: application/json

[
    {
        "issuer": "https://idp.example.org",
        "issuer_name#ja": "どっかの IdP",
        ...
    },
    ...
]
```

省略あり。


### 1.2. エラーレスポンス

エラーは [OAuth 2.0 Section 5.2] の形式で返す。



<!-- 参照 -->
[OAuth 2.0 Section 5.2]: http://tools.ietf.org/html/rfc6749#section-5.2
[OpenID Connect Discovery 1.0 Section 4.2]: http://openid.net/specs/openid-connect-discovery-1_0.html#ProviderConfigurationResponse
