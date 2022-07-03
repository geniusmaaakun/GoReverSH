# GoReverSH
Golang ReverseShell

##

## TODO
get hostname
get currentuser

color

clean

config

test

Refactoring

Makefile

暗号化
TLS -> AES共通鍵


その他
ゴルーチンリークを注意

dockerでテスト

テストコード

build op make　でmainの変数に代入

操作ログを残す
logfile

チャネルを使って、オブザーバーパターンで処理を分ける方がわかりやすい

killシグナルでクリーン終了

ダウンロードなどは、コネクション二つ貼ってやる方法もあり？
制御コマンド用と通信用

暗号化
base64でキーを埋め込み　RSA AESで暗号化
暗号文.共通鍵をサーバー側の秘密鍵で暗号化
クライアントは、公開鍵を受け取り複合化する

TLSを使う？
実際は証明書を発行しないので、独自に暗号化する必要がある