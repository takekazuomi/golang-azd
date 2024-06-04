# azd + go

## 目的

azd でazure infraをデプロイし、go でアプリを作成するというフローで開発する

もともと、azd で、Goがサポートされていれば良いが、当面無理そうなので、別途やることにする。

## azd インストール

```sh
$ curl -fsSL https://aka.ms/install-azd.sh | bash
azd-linux-amd64
NOTICE.txt
mkdir: cannot create directory ‘/opt/microsoft/azd’: Permission denied
install-azd: Creating /opt/microsoft/azd requires elevated permission. You may be prompted to enter credentials.
[sudo] password for takekazu:
install-azd: Writing to /opt/microsoft/azd/ requires elevated permission. You may be prompted to enter credentials.
install-azd: Writing to /usr/local/bin/ requires elevated permission. You may be prompted to enter credentials.
install-azd: Successfully installed to /opt/microsoft/azd
install-azd: Symlink created at /usr/local/bin/azd and pointing to /opt/microsoft/azd/azd-linux-amd64
install-azd:
install-azd: The Azure Developer CLI collects usage data and sends that usage data to Microsoft in order to help us improve your experience.
install-azd: You can opt-out of telemetry by setting the AZURE_DEV_COLLECT_TELEMETRY environment variable to 'no' in the shell you use.
install-azd:
install-azd: Read more about Azure Developer CLI telemetry: https://github.com/Azure/azure-dev#data-collection
$ azd version
azd version 1.9.3 (commit e1624330dcc7dde440ecc1eda06aac40e68aa0a3)
```

## azure infraをデプロイする

基本的にこれを使う
https://github.com/Azure-Samples/azd-starter-bicep


テンプレートを探すには、https://azure.github.io/awesome-azd/　を使うと便利

コマンドラインからは、`azd template list` で確認できる。

1. az initする
2. 適当にinfraの下をいじってデプロイ

## go のアプリ

下記の3つの下にある（もう1つ下げた方が良いかも？）

- cmd
- pkg
- api-specs

## TODO

- grpc reflection を入れる
- README更新
- storage access
- db access
- https call
- grpc call
- redis ?
- blog
