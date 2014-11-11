edo-idp-selector
===

ID プロバイダ一覧を選択するための UI。

### /

cookie の ID_PROVIDER_UUID に有効な ID プロバイダが設定されていれば、
クエリパラメータに id_provider_uuid と id_provider_login_uri を付加して、
/set_cookie にリダイレクト。
設定されていなければ、クエリパラメータを維持したまま /list にリダイレクト。

### /list

ID プロバイダ一覧を表示する。
各 ID プロバイダのリンクは、/set_cookie へのリンクに、
リクエストのクエリパラメータと id_provider_uuid を付加したもの。

### /set_cookie?id_provider_uuid={id_provide_uuid}

クエリパラメータに id_provider_uuid がなければ 400 Bad Request。
id_provider_uuid が有効な ID プロバイダでなければ 403 Forbidden。

id_provider_uuid の値を cookie の ID_PROVIDER_UUID に設定し、
クエリパラメータを維持したまま、id_provider_uuid のログイン URI にリダイレクト。
