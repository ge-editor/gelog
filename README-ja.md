# gelog

`gelog` は Go 標準の `log/slog` をベースにした、  
**シンプルで実用的なロガー初期化 & ハンドラ補助パッケージ**です。

- ログローテーション対応（lumberjack）
- 複数出力をまとめる `MultiHandler`
- 人間が読みやすい簡易フォーマット `PlainHandler`
- build tag による **debug / release 切り替え**

---

## 特徴

### ✅ slog ベース
Go 1.21 以降で標準化された `log/slog` をそのまま活かします。

### ✅ ログローテーション
`lumberjack` を利用し、ファイルサイズ・世代数・保存期間を管理。

### ✅ 複数出力対応
`MultiHandler` により、  
1つのログを **複数の slog.Handler に同時出力**できます。

### ✅ debug / release ビルド切り替え
`//go:build debug` により、  
- debug: ログ出力あり
- release: ログ呼び出し自体を無効化  

が可能です。

---

## インストール

```bash
go get github.com/ge-editor/gelog
```

## 使い方

```go
import "github.com/ge-editor/gelog"

func main() {
    gelog.InitLogger("app.log")

    gelog.Info("application started",
        "version", "1.0.0",
        "port", 8080,
    )
}
```

## ログ出力例

```
time=2025-01-01T12:00:00.000+09:00 level=INFO source=main.go:12 msg="application started" version=1.0.0 port=8080
```

## MultiHandler

複数の slog.Handler を束ねるためのハンドラです。

```go
handler := gelog.NewMultiHandler(
    fileHandler,
    stdoutHandler,
)
```

- Enabled / Handle / WithAttrs / WithGroup
- slog.Handler の仕様に準拠

## OS ごとの挙動

```go
if runtime.GOOS == "windows" {
    handler = fileHandler
} else {
    handler = NewMultiHandler(
        fileHandler,
        // Stdout は現在無効
    )
}
```

Windows: ファイル出力のみ

Unix 系: 複数出力を想定（現在はファイルのみ）

## debug / release ビルド

### debug ビルド

```bash
go build -tags debug
```

gelog.Info("this will be logged")

### release ビルド（デフォルト）

```bash
go build
```

gelog.Info("this will be ignored")

---

## VSCode 用 開発用環境変数の設定

`.vscode/settings.json`:

```json
{
    "go.toolsEnvVars": {
        "GOFLAGS": "-tags=debug,ja_JP"
    }
}
```

---
