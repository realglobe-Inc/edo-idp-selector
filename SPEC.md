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


## エンドポイント

|エンドポイント名|初期 URI|機能|
|:--|:--|:--|
|開始|/|選択処理を開始する|
|選択|/select|選択した IdP にリダイレクトさせる|
|UI|/html|UI 用のファイルを提供する|
|IdP 列挙|/api/list|UI 用に IdP 情報を提供する|


## 開始エンドポイント

処理を開始する。

リクエストに `prompt` パラメータを含み、その値が `select_account` を含む場合、セッションを選択中に設定しつつ、UI エンドポイントにリダイレクトさせる。
そうでなく、Cookie に X-Edo-Idp-Selector を含み、その値が有効なセッションの場合、セッションを更新しつつ、セッションに紐付く IdP にリダイレクトさせる。
そうでなければ、セッションを選択中に設定しつつ、UI エンドポイントにリダイレクトさせる。

|Cookie|値|
|:--|:--|
|X-Edo-Idp-Selector|セッション ID|

UI エンドポイントへのリダイレクト時には、選択チケットを発行し、セッションに紐付け、リダイレクト先のフラグメントにする。

```html
HTTP/1.1 302 Found
Location: /html#CgKa4ugl_k
```

## 選択エンドポイント

IdP の選択が終わった後の処理をする。

Cookie に X-Edo-Idp-Selector を含み選択中のセッションならば、リクエストから以下のパラメータを取り出す。

* **`idp`**
    * IdP の ID。
* **`ticket`**
    * 選択チケット。

選択チケットがセッションに紐付く選択チケットならば、セッションにこの IdP を紐付け、セッションを更新しつつ、IdP にリダイレクトさせる。


## UI エンドポイント

IdP 選択用の UI を提供する。
UI の目的は、以下のパラメータを選択エンドポイントに POST すること。

* **`idp`**
    * 選択した IdP の ID。
* **`ticket`**
    * 選択チケット。
      開始エンドポイントからのリダイレクト時にフラグメントとして付加される。


## IdP 列挙エンドポイント

UI 用の IdP 一覧を返す。

クエリで絞り込める。
クエリは以下の形式で連言で用いる。

```
<タグ名>=<該当する値の正規表現>
````

レスポンスは [OpenID Connect Discovery 1.0 Section 4.2] 形式の IdP 情報の JSON 配列である。

```json
[
    {
        "issuer": "https://idp.example.org",
        "friendly_name": "どっかの IdP",
        ...
    },
    ...
]
```


## エラーレスポンス

エラーは [OAuth 2.0 Section 4.1.2.1] の形式で返す。


<!-- 参照 -->
[OAuth 2.0 Section 4.1.2.1]: http://tools.ietf.org/html/rfc6749#section-4.1.2.1
[OpenID Connect Discovery 1.0 Section 4.2]: http://openid.net/specs/openid-connect-discovery-1_0.html#ProviderConfigurationResponse
