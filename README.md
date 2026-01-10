# Zuon (ズオン)

[![Release](https://img.shields.io/github/v/release/aomori446/zuon)](https://github.com/aomori446/zuon/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/aomori446/zuon)](https://goreportcard.com/report/github.com/aomori446/zuon)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Zuon** is a secure, cross-platform image steganography tool written in Go using the [Fyne](https://fyne.io/) toolkit. It allows you to hide sensitive text or files inside images using strong AES-256 encryption.

**Zuon**（ズオン）は、Go言語とFyneツールキットで構築された、セキュアなクロスプラットフォーム・ステガノグラフィツールです。強力なAES-256暗号化を使用して、画像の中に機密テキストやファイルを隠すことができます。

---

## 🚀 Features / 特徴

*   **Modern Sidebar UI**: A clean, professional interface for easy navigation.
*   **Secure Steganography**: Embeds data into the least significant bits (LSB) of the image, making it invisible to the naked eye.
*   **Unsplash Integration**: Search and use high-quality carrier images directly from the app.
*   **Strong Encryption**: All data is encrypted with **AES-GCM** using a user-provided password before embedding.
*   **Internationalization (i18n)**: Fully localized interface.
    *   🇺🇸 English
    *   🇨🇳 简体中文
    *   🇯🇵 日本語
*   **User-Friendly Error Handling**: Clear, localized error messages instead of technical jargon.

## 📦 Installation / インストール

### Download Binary (Recommended)
Go to the [Releases](https://github.com/aomori446/zuon/releases) page and download the latest version (`v1.3.0+`) for your operating system.

[リリースページ](https://github.com/aomori446/zuon/releases)から、お使いのOSに合わせた最新バージョン（`v1.3.0`以降）をダウンロードしてください。

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
2.  **Select Mode**: Use the sidebar to switch between "Embed" (隠す) and "Extract" (抽出).
3.  **Embed**:
    *   **Select Image**: Click the folder icon to open a local file, or click **"Search Web"** to find an image on Unsplash.
    *   **Input Data**: Enter your secret text or upload a file.
    *   **Set Password**: Set a strong password for encryption.
    *   **Save**: Generate and save the resulting image.
4.  **Extract**:
    *   Load the image containing hidden data.
    *   Enter the password used for encryption.
    *   Reveal the hidden message or save the extracted file.

### 🔑 Unsplash Configuration
To use the online image search feature, you will need a free **Unsplash Access Key**.
1.  Click the "Search Web" button in the app.
2.  A dialog will appear asking for your key.
3.  Enter your key, and it will be saved securely for future use.

(Developers can also set the `UNSPLASH_ACCESS_KEY` environment variable).

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

バグ報告や機能追加の提案など、プルリクエストは大歓迎です。

## 📄 License

MIT License - see the [LICENSE](LICENSE) file for details.
