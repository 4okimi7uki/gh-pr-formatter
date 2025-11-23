# gh-pr-formatter

![Go Version](https://img.shields.io/badge/Go-1.25-blue?logo=go)
![CI](https://github.com/4okimi7uki/gh-pr-formatter/actions/workflows/lint.yml/badge.svg)

GitHub CLI (`gh`) を使って **マージ済み Pull Request を自動収集 & Markdown に整形** する CLI ツールです。

- 直近のマージ PR を取得
- Author ごとに番号をグルーピング
- テンプレートを使って Markdown に整形

## 動作フロー

1. **環境チェック**

- `gh` がインストール済みか確認
- カレントディレクトリが `.git` 管理下か確認

2. **`main` ブランチの最新マージ日時を取得**

- `hotfix/xxx` 形式のブランチは除外して検索

3. **その日時〜現在までに `develop` にマージされた PR を収集**
4. **PR を `author（ログイン名）`ごとにグルーピング**
5. **Markdown へ整形し、以下の形式で出力**

```
./releasePrMarkdown/release_YYYYMMDD_hhmm.md
```

## CLI sample

```bash
========================================================
 Merged Pull Requests
========================================================
@4okimi7uki
- #18
- #16
- #15
- #13
...
========================================================
 🎉 SUCCESS: Release PR Markdown created successfully!
 -> Output: ./releasePrMarkdown/release_20251123_1918.md
========================================================
```

---

## 開発者向け

### Build（開発者向け）

Go のクロスコンパイル機能を利用して、各 OS（Linux / macOS / Windows ）向けの実行バイナリを生成できます。
成果物はすべて ./dist 配下に出力されます。

#### 個別 OS 用のバイナリを生成

環境に応じた 1 種類のバイナリを生成します。

```bash
make
```

#### すべての OS 向けバイナリをまとめて生成

```bash
make build-all
```

以下のような形式で出力されます：

```
dist/
├── gh-pr-formatter-mac-amd64
├── gh-pr-formatter-mac-arm64
├── gh-pr-formatter-linux-amd64
└── gh-pr-formatter-windows-amd64.exe
```

マルチプラットフォームに配布したいときは `make build-all` が便利です。

---

<small>2025 [Aoki Mizuki](https://github.com/4okimi7uki) – Developed with 🍭 and a sense of fun.</small>
