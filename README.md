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


# edo-idp-selector

IdP 一覧を選択するための UI バックエンド。


## 1. 起動


UI 用の HTML 等を ui ディレクトリの下に置く。

```
<任意のディレクトリ>/
├── edo-idp-selector
└── html
     ├── index.html
     ...
```

|オプション|値の意味・選択肢|
|:--|:--|
|-uiPath|UI 用 HTML 等を置くディレクトリパス|
|-uiUri|UI 用 HTML 等を提供する URI|


## 2. URI

|URI|機能|
|:--|:--|
|/list|IdP を列挙する API|
|/redirect|IdP を受け取り、リダイレクト処理をする|
|/html|UI 用の HTML を提供する|
|/|IdP の選択処理をする|

エラー時、リクエストに redirect_uri クエリが含まれている場合は、
OpenID Connect 式に redirect_uri へのリダイレクトに error クエリ等を付けてエラーを通知する。

### 2.1. GET /

prompt クエリが select_account の場合、クエリを維持したまま /html/index.html にリダイレクトする。

そうでなく、cookie の X-Edo-Idp-Id、または、idp クエリに有効な IdP が設定されている場合、
それを選択された IdP とみなして /redirect と同じ処理をする。
優先度は、

    idp クエリ > cookie の X-Edo-Idp-Id

そうでない場合、クエリを維持したまま /html/index.html にリダイレクトする。


### 2.2. GET /list

IdP 一覧を返す。
クエリで絞り込める。
クエリの形式は

    <タグ名>=<該当する値の正規表現>

レスポンスは JSON で、

```
[
    {
        "id"="https://example.com",
        "name"="どっかの IdP",
        ...
    },
    ...
]
```


### 2.3. GET /html/...

UI 用の HTML を提供する。

対応するディレクトリ内に、少なくとも index.html だけは置く必要がある。

UI の役目は、最終的に、/redirect に現在のクエリ及び選択した IdP を idp クエリとして付けた
以下のようなリンクを踏ませること。

    <a href="/redirect?client_id=...&idp=https%3A%2F%2Fexample.com">どっかの IdP</a>


### 2.4. GET /redirect

選択された IdP を受け取り、その認証 URI にクエリを維持したままリダイレクトさせる。
IdP の受け取りは、idp クエリ、または、idp フォームパラメータ。
レスポンス時に、Set-Cookie で X-Edo-Idp-Id を設定する。

/ を代わりに使うこともできるが、バグで無限ループしたら嫌なので、UI からは /redirect を踏ませる。


## 3. ライセンス

Apache License, Version 2.0
