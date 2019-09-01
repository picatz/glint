# glint

> ⚠️ Under development!

Extensible golang linting tool.

## Why?

After playing around a little bit with `go/ast` (and related packages), I realized there was an oppurtunity to create a more easily accessible golang linter using JSON-defined rules.

## Usage

Rules file:

```json
{
    "rules": [
        {
            "type": "",
            "comment": "",
            ...
        },
        ...
    ]
}
```

```console
$ glint examples/main.go
examples/main.go:4:2:we don't rely on these packages for almost anything
examples/main.go:8:2:we don't rely on these packages for almost anything
examples/main.go:11:2:we don't use golang.org/x/* packages
examples/main.go:15:58:use exactly 2048 for bits because examples are silly sometimes
examples/main.go:15:29:don't use math.Rand as source of entropy
examples/main.go:17:17:we don't use fmt.Errorf for some silly reaosn
examples/main.go:18:24:we don't use fmt.Errorf for some silly reaosn
examples/main.go:18:35:don't use uppercase error message string in fmt.Errorf formatted errors
examples/main.go:30:2:don't use http.Handle/http.HandleFunc which uses the DefaultServeMux due to possible side-effects
examples/main.go:31:2:don't use http.Handle/http.HandleFunc which uses the DefaultServeMux due to possible side-effects
```

> **Note**: Linting messages are output to STDOUT in the following format `file:line:column:comment`

## Rule Types

All rules require a `type` and a `comment` to define what it is and what message is included in the linting error. For each supported rule type, there are a collection of options which help customize how the rule should inspect the program. This attempts to be as declarative as possible.

```json
{
    "type": "",
    "comment": "",
    ...
}
```

### Import

Using the `"import"` type you can specify what imported packages should or should not be included in a program.

#### Available Options for Import

| Option        | Description                                             | Required  |
| ------------- |:--------------------------------------------------------|----------:|
| `cannot_match`| array of regular expressions for packages the program **should not** use | false     |
| `must_match`  | array of regular expressions the packages program **should** use     | false     |

```json
{
    "type": "import",
    "comment": "don't use golang.org/x/* packages",
    "cannot_match": [
        "golang.org/x/\\w+"
    ]
}
```

### Method

Using the `"method"` type you can define certain method calls that should not be used.

#### Available Options for Method

| Option           | Description                                                                    | Required  |
| ---------------- |:-------------------------------------------------------------------------------|----------:|
| `call`           | the package.Method call to inspect                                             | false     |
| `argument`       | the index of the method argument (starting at `0` for the first) to inspect    | false     |
| `less_than`      | the method argument, if it's an `int`, must be less than the given value       | false     |
| `greater_than`   | the method argument, if it's an `int`, must be greater than the given value    | false     |
| `equals`         | the method argument, if it's an `int`, must equal exactly the given value      | false     |
| `dont_use`       | the method call should not be used if the given faluse is true                 | false     |
| `cannot_match`   | array of regular expressions for method calls or arguments *should not** match | false     |

```json
{
   "type": "method",
   "comment": "use EXACTLY 2048 bits when generating RSA keys for some reason",
   "call": "rsa.GenerateKey",
   "argument": 1,
   "less_than": 2049,
   "greater_than": 0,
   "equals": 2048
}
```

```json
{
    "type": "method",
    "comment": "we don't use fmt.Errorf for some reaosn",
    "call": "fmt.Errorf",
    "dont_use": true
}
```

```json
{
    "type": "method",
    "comment": "don't use uppercase error message string in fmt.Errorf formatted errors for some reason",
    "call": "fmt.Errorf",
    "argument": 0,
    "cannot_match": [
        "^[A-Z]+"
    ]
}
```

```json
{
    "type": "method",
    "comment": "don't use math.Rand as source of entropy",
    "cannot_match": [
        "rand.New$"
    ]
}
```

```json
{
    "type": "method",
    "comment": "don't use http.Handle/http.HandleFunc which uses the DefaultServeMux due to possible side-effects",
    "cannot_match": [
                "http.Handle$",
                "http.HandleFunc$"
    ]
}
```
