# gelog

`gelog` is a **simple and practical logger initialization and handler helper package**
built on top of Go’s standard `log/slog`.

- Log rotation support (via lumberjack)
- `MultiHandler` for fan-out logging
- Human-readable simple format (`PlainHandler`)
- **debug / release switching** using build tags

---

## Features

### ✅ Based on slog
Uses Go’s standard `log/slog`, officially available since Go 1.21.

### ✅ Log rotation
Uses `lumberjack` to manage log file size, backups, and retention period.

### ✅ Multiple outputs
With `MultiHandler`, a single log record can be sent to **multiple slog.Handlers**.

### ✅ debug / release build switching
Using `//go:build debug`:

- **debug**: logging enabled
- **release**: logging calls become no-ops

This allows zero logging overhead in production builds.

---

## Installation

```bash
go get github.com/ge-editor/gelog
```

---

## Usage

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

---

## Log Output Example

```
time=2025-01-01T12:00:00.000+09:00 level=INFO source=main.go:12 msg="application started" version=1.0.0 port=8080
```

---

## MultiHandler

`MultiHandler` bundles multiple `slog.Handler` instances into one.

```go
handler := gelog.NewMultiHandler(
    fileHandler,
    stdoutHandler,
)
```

- Implements `Enabled`, `Handle`, `WithAttrs`, `WithGroup`
- Fully compatible with `slog.Handler`

---

## OS-specific Behavior

```go
if runtime.GOOS == "windows" {
    handler = fileHandler
} else {
    handler = NewMultiHandler(
        fileHandler,
        // Stdout output is currently disabled
    )
}
```

- **Windows**: file output only
- **Unix-like systems**: designed for multiple outputs (currently file only)

---

## debug / release Builds

### debug build

```bash
go build -tags debug
```

```go
gelog.Info("this will be logged")
```

### release build (default)

```bash
go build
```

```go
gelog.Info("this will be ignored")
```

---

## VSCode Development Environment Variables

`.vscode/settings.json`:

```json
{
    "go.toolsEnvVars": {
        "GOFLAGS": "-tags=debug,ja_JP"
    }
}
```

---
