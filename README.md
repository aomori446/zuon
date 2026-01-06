# zuon (ズオン)

Go言語で書かれたシンプルなステガノグラフィツールです。画像ファイルの中に暗号化されたデータを隠すことができます。

## 特徴

*   **AES-GCM暗号化**: データはパスワードを使って強力に暗号化されます。
*   **LSBステガノグラフィ**: 画像のピクセルデータの下位ビットに情報を埋め込みます。
*   **NRGBA対応**: 画像は内部的にNRGBA形式で処理されます。

## 使い方

このライブラリは `internal` パッケージとして実装されています。

### データの埋め込み (Embed)

```go
// 画像、隠したいデータ、オフセット、パスワードを指定して埋め込みます
resultImage, err := internal.EmbedData(srcImage, []byte("秘密のメッセージ"), 0, "my-password")
```

### データの抽出 (Extract)

```go
// 画像、オフセット、パスワードを指定してデータを取り出します
data, err := internal.ExtractData(encodedImage, 0, "my-password")
```

## インストール

```bash
go get github.com/aomori446/zuon
```

## ライセンス

MIT License
