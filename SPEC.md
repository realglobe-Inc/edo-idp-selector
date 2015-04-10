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

以降の動作記述において、箇条書きに以下の構造を持たせることがある。

* if
    * then
* else if
    * then
* else


## 1. エンドポイント

|エンドポイント名|初期パス|機能|
|:--|:--|:--|
|開始|/|選択処理を開始する|
|選択|/select|選択した IdP にリダイレクトさせる|
|UI|/ui/index.html|UI を提供する|
|IdP 列挙|/issinfo|UI 用に IdP 情報を提供する|


## 2. セッション

開始および選択エンドポイントではセッションを利用する。

|Cookie 名|値|
|:--|:--|
|Idp-Selector|セッション ID|

開始および選択エンドポイントへのリクエスト時に、セッション ID が通知されなかった場合、セッションを発行する。
セッションの期限に余裕がない場合、設定を引き継いだセッションを発行する。

開始および選択エンドポイントからのレスポンス時に、未通知のセッション ID を通知する。


## 3. 開始エンドポイント

処理を開始する。

まず、セッションにリクエスト内容を紐付ける。

* リクエストに `prompt` パラメータを含み、その値が `select_account` を含む場合、
    * UI エンドポイントにリダイレクトさせる。
* そうでなく、IdP に紐付くセッションである場合、
    * IdP にリダイレクトさせる。
* そうでなければ、UI エンドポイントにリダイレクトさせる。

UI エンドポイントへのリダイレクト時には、チケットを発行する。
チケットをセッションに紐付ける。
チケットをフラグメントとして付加した UI エンドポイントにリダイレクトさせる。

IdP へのリダイレクト時には、セッションとリクエスト内容やチケットとの紐付けを解く。
リクエスト内容を付加した IdP のユーザー認証エンドポイントにリダイレクトさせる。


### 3.1. リクエスト例

```http
GET /?response_type=code%20id_token&scope=openid
    &client_id=https%3A%2F%2Fta.example.org
    &redirect_uri=https%3A%2F%2Fta.example.org%2Freturn&state=Ito-lCrO2H
    &nonce=v46QjbP6Qr HTTP/1.1
Host: selector.example.org
Cookie: Idp-Selector=caiQ2D0ab04N0EPdCcG2OnB4SyBnME
```

改行とインデントは表示の都合による。


### 3.2. レスポンス例

UI へのリダイレクト例。

```http
HTTP/1.1 302 Found
Location: /ui/index.html#CgKa4ugl_k
```

改行とインデントは表示の都合による。


## 4. 選択エンドポイント

IdP が選択された後の処理をする。

* チケットと紐付くセッションでない場合、
    * エラーを返す。
* そうでなければ、リクエストから以下のパラメータを取り出す。

|パラメータ名|必要性|値|
|:--|:--|:--|
|**`ticket`**|必須|チケット|
|**`issuer`**|必須|選択された IdP の ID|
|**`locale`**|任意|選択された表示言語|

* チケットがセッションに紐付くものと異なる、または、IdP が正当でない場合、
    * エラーを返す。
* そうでなければ、設定を引き継いだセッションを発行する。
  IdP をセッションに紐付ける。
  IdP にリダイレクトさせる。


### 4.1. リクエスト例

```http
POST /select HTTP/1.1
Host: selector.example.org
Cookie: Idp-Selector=caiQ2D0ab04N0EPdCcG2OnB4SyBnME
Content-Type: application/x-www-form-urlencoded

ticket=CgKa4ugl_k&issuer=https%3A%2F%2Fidp.example.org
```


### 4.2. レスポンス例

```http
HTTP/1.1 302 Found
Set-Cookie: Idp-Selector=gWWw7dOxT0Op3bPV6vUHGr16hrg0Q4;
    Expires=Tue, 24 Mar 2015 01:59:23 GMT; Path=/; Secure; HttpOnly
Location: https://idp.example.org/auth?response_type=code%20id_token
    &scope=openid&client_id=https%3A%2F%2Fta.example.org
    &redirect_uri=https%3A%2F%2Fta.example.org%2Freturn&state=Ito-lCrO2H
    &nonce=v46QjbP6Qr
```

改行とインデントは表示の都合による。


## 5. UI エンドポイント

IdP 選択用の UI を提供する。

以下のパラメータを受け付ける。

|パラメータ名|必要性|値|
|:--|:--|:--|
|**`issuers`**|任意|特に候補になる IdP の ID の JSON 配列|
|**`display`**|任意|[OpenID Connect Core 1.0 Section 3.1.2.1] の `display` と同じもの|
|**`locales`**|任意|[OpenID Connect Core 1.0 Section 3.1.2.1] の `ui_locales` と同じもの|

UI の目的は、選択エンドポイントに POST させること。


### 5.1. リクエスト例

```http
GET /ui/index.html HTTP/1.1
Host: selector.example.org
```


## 6. IdP 列挙エンドポイント

UI 用に IdP 一覧を返す。

クエリで絞り込める。
クエリは以下の形式の連言として用いる。

```
<タグ名>=<該当する値の正規表現>
```

レスポンスは [OpenID Connect Discovery 1.0 Section 4.2] 形式の IdP 情報の JSON 配列である。
ただし、以下の最上位要素を加える。

* **`friendly_name`**
    * 名前。
      言語タグが付くことがある。


### 6.1. リクエスト例

```http
GET /issinfo?issuer=%5C.example%5C.org%24
Host: selector.example.org
```


### 6.2. レスポンス例

```http
HTTP/1.1 200 OK
Content-Type: application/json

[
    {
        "issuer": "https://idp.example.org",
        "friendly_name#ja": "どっかの IdP",
        ...
    },
    ...
]
```

省略あり。


## 7. エラーレスポンス

エラーは [OAuth 2.0 Section 4.1.2.1] の形式で返す。

セッションがある場合、セッションとリクエスト内容やチケットとの紐付けを解く。


## 8. 外部データ

以下に分ける。

* 共有データ
    * 他のプログラムと共有する可能性のあるもの。
* 非共有データ
    * 共有するとしてもこのプログラムの別プロセスのみのもの。


### 8.1. 共有データ


#### 8.1.1. IdP 情報

以下を含む。

* ID
* 名前
* 認証エンドポイント

以下の操作が必要。

* ID による取得
* 全取得
    * 任意のタグに対する正規表現による絞り込み


#### 8.1.2. TA 情報

エラーレスポンス時に必要。
以下を含む。

* ID
* リダイレクト URI

以下の操作が必要。

* ID による取得


### 8.2. 非共有データ


#### 8.2.1. セッション

以下を含む。

* ID \*
* 有効期限 \*
* 選択した IdP の ID \*
* リクエスト内容
* チケット
* 過去に選択した IdP の ID
* UI 表示言語

\* は設定を引き継がない。

以下の操作が必要。

* 保存
* ID による取得
* 上書き
    * ID、有効期限以外。


<!-- 参照 -->
[OAuth 2.0 Section 4.1.2.1]: http://tools.ietf.org/html/rfc6749#section-4.1.2.1
[OpenID Connect Core 1.0 Section 3.1.2.1]: http://openid-foundation-japan.github.io/openid-connect-core-1_0.ja.html#AuthRequest
[OpenID Connect Discovery 1.0 Section 4.2]: http://openid.net/specs/openid-connect-discovery-1_0.html#ProviderConfigurationResponse
