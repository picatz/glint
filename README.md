# glint

Extensible golang linting tool.

## Usage

```console
$ glint examples/main.go
examples/main.go:4:2:we don't rely on these packages for almost anything
examples/main.go:8:2:we don't rely on these packages for almost anything
examples/main.go:10:2:we don't use golang.org/x/* packages
examples/main.go:14:58:use exactly 2048 for bits because examples are silly sometimes
examples/main.go:14:29:don't use math.Rand as source of entropy
examples/main.go:16:17:we don't use fmt.Errorf for some silly reaosn
examples/main.go:17:24:we don't use fmt.Errorf for some silly reaosn
examples/main.go:17:35:don't use uppercase error message string in fmt.Errorf formatted errors
```
