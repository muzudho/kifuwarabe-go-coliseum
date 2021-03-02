# Set up app on Windows

## Pre-install

```console
go mod init
```

## Install

Build:  

```shell
# 使っていないパッケージを、インストールのリストから削除するなら
# go mod tidy

# 自作のパッケージを更新(再インストール)したいなら
# go get -u all

go build
```
