# Zuon (ズオン)

[![Release](https://img.shields.io/github/v/release/aomori446/zuon)](https://github.com/aomori446/zuon/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/aomori446/zuon)](https://goreportcard.com/report/github.com/aomori446/zuon)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Zuon** is a secure, cross-platform image steganography tool written in Go using the [Fyne](https://fyne.io/) toolkit. It allows you to hide sensitive text or files inside images using strong AES-256 encryption.

**Zuon**（ズオン）は、Go言語とFyneツールキットで構築された、セキュアなクロスプラットフォーム・ステガノグラフィツールです。強力なAES-256暗号化を使用して、画像の中に機密テキストやファイルを隠すことができます。

---

## 🚀 Features / 特徴

*   **Secure Steganography**: Embeds data into the least significant bits (LSB) of the image, making it invisible to the naked eye.
*   **Strong Encryption**: All data is encrypted with **AES-GCM** using a user-provided password before embedding.
*   **Cross-Platform UI**: Modern, easy-to-use graphical interface (Windows, macOS, Linux).
*   **Internationalization (i18n)**: Fully localized interface.
    *   🇺🇸 English
    *   🇯🇵 日本語
    *   🇨🇳 简体中文
*   **Responsive**: Heavy encryption tasks run in the background, keeping the UI smooth.

## 📦 Installation / インストール

### Download Binary (Recommended)
Go to the [Releases](https://github.com/aomori446/zuon/releases) page and download the latest version for your operating system.

[リリースページ](https://github.com/aomori446/zuon/releases)から、お使いのOSに合わせた最新バージョンをダウンロードしてください。

### Build from Source
If you have Go installed, you can build it yourself:

```bash
# Clone the repository
git clone https://github.com/aomori446/zuon.git
cd zuon

# Install dependencies
go mod tidy

# Run the app
go run cmd/zuon/main.go
```

## 🛠 Usage / 使い方

1.  **Launch Zuon**: Open the application.
2.  **Select Mode**: Choose between "Embed" (隠す) or "Extract" (抽出).
3.  **Embed**:
    *   Load a carrier image.
    *   Enter your secret text or upload a file.
    *   Set a strong password.
    *   Save the resulting image.
4.  **Extract**:
    *   Load the image containing hidden data.
    *   Enter the password used for encryption.
    *   Reveal the hidden message or save the extracted file.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

バグ報告や機能追加の提案など、プルリクエストは大歓迎です。

## 📄 License

MIT License - see the [LICENSE](LICENSE) file for details.