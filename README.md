# Polyglot Phrase Lab

中国語の例文を登録し、語順練習を行う学習アプリの公開用プロトタイプです。

このリポジトリは、12言語を別々のサンプルとして並べるのではなく、ひとつのアプリを構成するコンポーネントとしてまとめています。共通のデータ形式とHTTP APIを境界にして、Web・iOS・Android・管理・学習ロジック・音声処理を接続します。

## 構成と言語の役割

| 言語 | コンポーネント | 役割 |
|---|---|---|
| Go | `services/api` | 例文・並び替え練習のHTTP API |
| Rust | `engines/order` | 語句順の採点コア |
| C / C++ | `engines/audio` | 音声の無音区間・音量特徴の抽出 |
| JavaScript | `apps/web` | ブラウザー学習画面 |
| TypeScript | `packages/contracts` | APIデータ契約と型付きクライアント |
| Swift | `apps/ios` | iOSクライアントのAPI連携例 |
| Kotlin | `apps/android` | AndroidクライアントのAPI連携例 |
| Java | `services/review` | 間隔反復による次回復習日の計算 |
| Ruby | `tools/content-importer` | CSV教材を共通JSONへ変換 |
| PHP | `apps/admin` | 教材管理用の小さな管理画面 |
| C# | `tools/content-checker` | 教材JSONの品質検査 |

## 実装状況

- Webの並び替え練習とGo APIは動作する構成です。
- Rust採点エンジンはGo APIから実行できます。
- Ruby教材取込、C#教材検査、Java復習日計算、C/C++音声特徴抽出は独立して実行できる補助機能です。
- Swift / KotlinはAPI連携コードです。ネイティブアプリの画面・署名・配布設定は含みません。
- 音声解析・復習計算はWeb画面/APIにはまだ接続していません。発音の正誤を判定する完成機能ではありません。

## データ契約

例文は `data/phrases.json` に保存します。各レコードは `id`, `level`, `zh`, `pinyin`, `ja`, `tokens` を持ちます。API仕様は `contracts/openapi.yaml` を基準とします。

## 起動

Go 1.22以降とRust toolchainがある環境で:

```sh
cd services/api
cargo build --release --manifest-path ../../engines/order/Cargo.toml
PHRASE_ORDER_BIN="../../engines/order/target/release/phrase-order" go run .
```

ブラウザー画面は別ターミナルで:

```sh
cd apps/web
node server.mjs
```

`http://localhost:4173` を開きます。APIは `http://localhost:8080` です。Rustを使わずGo内蔵の比較で動かす場合は、`PHRASE_ORDER_BIN` を省略できます。

## 検証

```sh
cd services/api && go test ./...
cd ../../engines/audio && cmake -S . -B build && cmake --build build && ctest --test-dir build
```

各言語コンポーネントの詳細と実装状況は `docs/architecture.md` を参照してください。

## スコアについて

このコードがFindyの特定の偏差値を保証するものではありません。公開する場合は、履歴や設定ファイル、教材の権利、秘密情報を確認してください。
