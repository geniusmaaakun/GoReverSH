# GoReverSH
![ソースコードサイズ](https://img.shields.io/github/repo-size/geniusmaaakun/GoReverSH)

![ダウンロード]https://img.shields.io/github/downloads/geniusmaaakun/GoReverSH/total

![ライセンス]https://img.shields.io/apm/l/GoReverSH

![バージョン]https://img.shields.io/amo/v/GoReverSH

![レーティング]https://img.shields.io/amo/stars/GoReverSH

Golang ReverseShell
other en@ble@ny

# 詳細

# 使い方



## 未実装
### Test
TestWaitNotice
TestClist
Cert

### 関数
get hostname
get currentuser


## TODO
makefile
変数　ビルドファイル名など

ReadMe

color

clean
    ファイル全削除


upload　土日
    dir対応
    downloadを逆にする。
    クライアント側もオブザーバーパターンにする方がいい
    outをデコード　

upload
特定のUrlのzipダウンロード、解答するコードを送りつけて実行するプログラム開発
もしくはDirectoryの場合は、zipを送りつけて解凍


ダウンロード
    サイズ、パーミッション

出力ディレクトリ変更コマンド


別プログラムに偽装　バッチファイルなど
管理権限


暗号化
TLS -> AES共通鍵


クライアントテストコード 
計画して一つずつ
connモックを使って読み込みするRead


ロゴとアスキーアート


## その他
ゴルーチンリークを注意

build op make　でmainの変数に代入

操作ログを残す
logfile
出力をsetOutputで書き出し


ダウンロードなどは、コネクション二つ貼ってやる方法もあり？
制御コマンド用と通信用


暗号化
base64でキーを埋め込み　RSA AESで暗号化
暗号文.共通鍵をサーバー側の秘密鍵で暗号化
クライアントは、公開鍵を受け取り複合化する


TLSを使う？
実際は証明書を発行しないので、独自に暗号化する必要がある