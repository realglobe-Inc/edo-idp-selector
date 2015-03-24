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


# edo-idp-selector の仕様（目標）

ユーザーに IdP を選択させて、その IdP にユーザーを受け渡す。


## 1. エンドポイント

|エンドポイント名|初期 URI|機能|
|:--|:--|:--|
|開始|/|選択処理を開始する|
|選択|/select|選択した IdP にリダイレクトさせる|
|UI|/html/index.html|UI を提供する|
|IdP 列挙|/idpinfo|UI 用に IdP 情報を提供する|


### 1.1. 記法

以降の動作仕様の記述において、箇条書きに以下の構造を持たせることがある。

* if
    * then
* else if
    * then
* else


## 2. 開始エンドポイント

処理を開始する。

* リクエストに `prompt` パラメータを含み、その値が `select_account` を含む場合、
    * UI エンドポイントにリダイレクトさせる。
* そうでなく、Cookie に X-Edo-Idp-Selector を含み、その値が有効なセッションの場合、
    * セッションに紐付く IdP にリダイレクトさせる。
* そうでなければ、UI エンドポイントにリダイレクトさせる。

|Cookie ラベル|値|
|:--|:--|
|X-Edo-Idp-Selector|セッション ID|

UI エンドポイントへのリダイレクト時には、選択チケットを発行し、それをセッションに紐付け、セッションを更新しつつ、選択チケットをフラグメントとして付加した UI エンドポイントにリダイレクトさせる。

IdP へのリダイレクト時には、IdP をセッションに紐付け、セッションを更新しつつ、リダイレクトさせる。


### 2.1. リクエスト例

```http
GET /?response_type=code%20id_token&scope=openid
    &client_id=https%3A%2F%2Fta.example.org
    &redirect_uri=https%3A%2F%2Fta.example.org%2Freturn
    &state=Ito-lCrO2H&nonce=v46QjbP6Qr HTTP/1.1
Host: selector.example.org
```

改行とインデントは表示の都合による。


### 2.2. レスポンス例

UI へのリダイレクト例。

```http
HTTP/1.1 302 Found
Set-Cookie: X-Edo-Idp-Selector=caiQ2D0ab04N0EPdCcG2OnB4SyBnMEfmfoEiv6JF;
    Path=/; Expires=Tue, 24 Mar 2015 01:59:21 GMT
Location: /html/index.html#CgKa4ugl_k
```

改行とインデントは表示の都合による。


## 3. 選択エンドポイント

IdP が選択された後の処理をする。

* Cookie に X-Edo-Idp-Selector を含まない、または、選択チケットと紐付くセッションでない場合、
    * エラーを返す。
* そうでなければ、リクエストから以下のパラメータを取り出す。

|パラメータ名|必要性|値|
|:--|:--|:--|
|**`ticket`**|必須|選択チケット|
|**`idp`**|必須|選択された IdP の ID|

* 選択チケットがセッションに紐付くものと異なる、または、IdP が正当でない場合、
    * エラーを返す。
* そうでなければ、セッションを更新しつつ、IdP にリダイレクトさせる。


### 3.1. リクエスト例

```http
POST /select HTTP/1.1
Host: selector.example.org
Cookie: X-Edo-Idp-Selector=caiQ2D0ab04N0EPdCcG2OnB4SyBnMEfmfoEiv6JF
Content-Type: application/x-www-form-urlencoded

ticket=CgKa4ugl_k&idp=https%3A%2F%2Fidp.example.org
```


### 3.2. レスポンス例

```http
HTTP/1.1 302 Found
Set-Cookie: X-Edo-Idp-Selector=caiQ2D0ab04N0EPdCcG2OnB4SyBnMEfmfoEiv6JF;
    Path=/; Expires=Tue, 24 Mar 2015 01:59:23 GMT
Location: https://idp.example.org/auth?
    response_type=code%20id_token&scope=openid
    &client_id=https%3A%2F%2Fta.example.org
    &redirect_uri=https%3A%2F%2Fta.example.org%2Freturn
    &state=Ito-lCrO2H&nonce=v46QjbP6Qr
```

改行とインデントは表示の都合による。


## 4. UI エンドポイント

IdP 選択用の UI を提供する。

以下のパラメータを受け付ける。

|パラメータ名|必要性|値|
|:--|:--|:--|
|**`idps`**|任意|特に候補になる IdP の ID の JSON 配列|
|**`display`**|任意|[OpenID Connect Core 1.0 Section 3.1.2.1] の `display` と同じもの|
|**`locales`**|任意|[OpenID Connect Core 1.0 Section 3.1.2.1] の `ui_locales` と同じもの|

UI の目的は、選択エンドポイントに POST させること。


### 4.1. リクエスト例

```http
GET /html/index.html HTTP/1.1
Host: selector.example.org
Cookie: X-Edo-Idp-Selector=caiQ2D0ab04N0EPdCcG2OnB4SyBnMEfmfoEiv6JF
```


## 5. IdP 列挙エンドポイント

UI 用の IdP 一覧を返す。

クエリで絞り込める。
クエリは以下の形式の連言として用いる。

```
<タグ名>=<該当する値の正規表現>
````

レスポンスは [OpenID Connect Discovery 1.0 Section 4.2] 形式の IdP 情報の JSON 配列である。


### 5.1. リクエスト例

```http
GET /idpinfo?issuer=%5C.example%5C.org%24
Host: selector.example.org
Cookie: X-Edo-Idp-Selector=caiQ2D0ab04N0EPdCcG2OnB4SyBnMEfmfoEiv6JF
```


### 5.2. レスポンス例

```http
HTTP/1.1 200 OK
Content-Type: application/json

[
    {
        "issuer": "https://idp.example.org",
        "friendly_name": "どっかの IdP",
        ...
    },
    ...
]
```

省略あり。


## 6. エラーレスポンス

エラーは [OAuth 2.0 Section 4.1.2.1] の形式で返す。


## 7. 外部データ

以下に分ける。

* 共有データ
    * 他のプログラムと共有する可能性のあるもの。
* 非共有データ
    * 共有するとしてもこのプログラムの別プロセスのみのもの。


### 7.1. 共有データ

* IdP 情報
    * ID
    * 認証エンドポイント


### 7.2. 非共有データ

* セッション
    * ID
    * 有効 / 無効
    * 過去に選択した IdP の ID
    * 現在のリクエストパラメータ
    * 選択チケット


<!-- 参照 -->
[OAuth 2.0 Section 4.1.2.1]: http://tools.ietf.org/html/rfc6749#section-4.1.2.1
[OpenID Connect Core 1.0 Section 3.1.2.1]: http://openid-foundation-japan.github.io/openid-connect-core-1_0.ja.html#AuthRequest
[OpenID Connect Discovery 1.0 Section 4.2]: http://openid.net/specs/openid-connect-discovery-1_0.html#ProviderConfigurationResponse
