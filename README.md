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

ID プロバイダ選択サービス。


## 1. インストール

[go] が必要。
go のインストールは http://golang.org/doc/install を参照のこと。

go をインストールしたら、

```shell
go get github.com/realglobe-Inc/edo-idp-selector
```

適宜、依存ライブラリを `go get` すること。


## 2. 実行

以下ではバイナリファイルが `${GOPATH}/bin/edo-idp-selector` にあるとする。
パスが異なる場合は適宜置き換えること。


### 2.1. UI の準備

選択 UI を edo-idp-selector で提供する場合は、適当なディレクトリに UI 用ファイルを用意する。

```
<UI ディレクトリ>/
├── select.html
...
```

UI ディレクトリは起動オプションで指定する。


### 2.2. 起動

単独で実行できる。

```shell
${GOPATH}/bin/edo-idp-selector
```

### 2.3. 起動オプション

|オプション名|初期値|値|
|:--|:--|:--|
|-uiDir||UI 用ファイルを置くディレクトリパス|


### 2.4. デーモン化

単独ではデーモンとして実行できないため、[Supervisor] 等と組み合わせて行う。


## 3. 動作仕様

ユーザーに IdP を選択させて、その IdP にユーザーを受け渡す。

### 3.1. エンドポイント

|エンドポイント名|初期パス|機能|
|:--|:--|:--|
|開始|/|[ID プロバイダ選択機能](/page/idpselect)を参照|
|選択|/select|[ID プロバイダ選択機能](/page/idpselect)を参照|
|選択 UI|/ui/select.html|[ID プロバイダ選択機能](/page/idpselect)を参照|
|ID プロバイダ列挙|/api/info/issuer|[ID プロバイダ情報提供機能](/api/idp)を参照|


## 4. API

[GoDoc](http://godoc.org/github.com/realglobe-Inc/edo-idp-selector)


## 5. ライセンス

Apache License, Version 2.0

[Supervisor]: http://supervisord.org/
[go]: http://golang.org/
