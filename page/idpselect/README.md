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


# ID プロバイダ選択機能

ユーザーに ID プロバイダを選択させて、その ID プロバイダにユーザーを受け渡す。


## 1. 動作仕様

以降、箇条書きに以下の構造を持たせることがある。

* if
    * then
* else if
    * then
* else


### 1.1. エンドポイント

|エンドポイント名|機能|
|:--|:--|
|開始|選択処理を開始する|
|選択|選択した ID プロバイダにリダイレクトさせる|
|選択 UI|選択 UI を提供する|


### 1.2. セッション

開始、選択エンドポイントではセッションを利用する。

開始、選択エンドポイントへのリクエスト時に、有効なセッションが宣言されなかった場合、セッションを発行する。
開始エンドポイントへのリクエスト時に、セッションの期限に余裕がない場合、設定を引き継いだセッションを発行する。

開始、選択エンドポイントからのレスポンス時に、未通知のセッション ID を通知する。


### 1.3. 開始エンドポイント

ID プロバイダ選択処理を開始する。

まず、セッションにリクエスト内容を紐付ける。

* リクエストに `prompt` パラメータを含み、その値が `select_account` を含む場合、
    * 選択 UI エンドポイントにリダイレクトさせる。
* そうでなく、ID プロバイダに紐付くセッションである場合、
    * ID プロバイダにリダイレクトさせる。
* そうでなければ、選択 UI エンドポイントにリダイレクトさせる。

選択 UI エンドポイントへのリダイレクト時には、チケットを発行する。
チケットをセッションに紐付ける。
チケットをフラグメントとして付加した選択 UI エンドポイントにリダイレクトさせる。

ID プロバイダへのリダイレクト時には、セッションからリクエスト内容とチケットへの紐付けを解く。
リクエスト内容を付加した ID プロバイダのユーザー認証エンドポイントにリダイレクトさせる。


#### 1.3.1. リクエスト例

```http
GET /?response_type=code%20id_token&scope=openid
    &client_id=https%3A%2F%2Fta.example.org
    &redirect_uri=https%3A%2F%2Fta.example.org%2Freturn&state=Ito-lCrO2H
    &nonce=v46QjbP6Qr HTTP/1.1
Host: selector.example.org
Cookie: Idp-Selector=caiQ2D0ab04N0EPdCcG2OnB4SyBnME
```

改行とインデントは表示の都合による。


#### 1.3.2. レスポンス例

選択 UI へのリダイレクト例。

```http
HTTP/1.1 302 Found
Location: /ui/select.html#CgKa4ugl_k
```

改行とインデントは表示の都合による。


### 1.4. 選択エンドポイント

ID プロバイダが選択された後の処理をする。

* チケットと紐付くセッションでない場合、
    * エラーを返す。
* そうでなければ、リクエストから以下のパラメータを取り出す。

|パラメータ名|必要性|値|
|:--|:--|:--|
|**`ticket`**|必須|チケット|
|**`issuer`**|必須|選択された ID プロバイダの ID|
|**`locale`**|任意|選択された表示言語|

* チケットがセッションに紐付くものと異なる、または、ID プロバイダが正当でない場合、
    * エラーを返す。
* そうでなければ、ID プロバイダをセッションに紐付ける。
  ID プロバイダにリダイレクトさせる。


#### 1.4.1. リクエスト例

```http
POST /select HTTP/1.1
Host: selector.example.org
Cookie: Idp-Selector=caiQ2D0ab04N0EPdCcG2OnB4SyBnME
Content-Type: application/x-www-form-urlencoded

ticket=CgKa4ugl_k&issuer=https%3A%2F%2Fidp.example.org
```


#### 1.4.2. レスポンス例

```http
HTTP/1.1 302 Found
Location: https://idp.example.org/auth?response_type=code%20id_token
    &scope=openid&client_id=https%3A%2F%2Fta.example.org
    &redirect_uri=https%3A%2F%2Fta.example.org%2Freturn&state=Ito-lCrO2H
    &nonce=v46QjbP6Qr
```

改行とインデントは表示の都合による。


### 1.5. 選択 UI エンドポイント

本パッケージでは提供しない。
以下は要件。

ID プロバイダ選択用の UI を提供する。
UI の目的は、選択エンドポイントに POST させること。

以下のパラメータを受け付ける。

|パラメータ名|必要性|値|
|:--|:--|:--|
|**`issuers`**|任意|特に候補になる ID プロバイダの ID の JSON 配列|
|**`display`**|任意|[OpenID Connect Core 1.0 Section 3.1.2.1] の `display` と同じもの|
|**`locales`**|任意|[OpenID Connect Core 1.0 Section 3.1.2.1] の `ui_locales` と同じもの|


## 2. エラーレスポンス

エラーは [OAuth 2.0 Section 4.1.2.1] の形式で返す。
セッションからリクエスト内容とチケットへの紐付けを解く。


<!-- 参照 -->
[OAuth 2.0 Section 4.1.2.1]: http://tools.ietf.org/html/rfc6749#section-4.1.2.1
[OpenID Connect Core 1.0 Section 3.1.2.1]: http://openid-foundation-japan.github.io/openid-connect-core-1_0.ja.html#AuthRequest
