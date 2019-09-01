# glint

> ⚠️ Under development!

Extensible golang linting tool.

## Install

```console
$ go get -u github.com/picatz/glint
```

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
examples/main.go:9:2:we don't rely on these packages for almost anything
examples/main.go:13:2:don't use the unsafe package
examples/main.go:15:2:we don't use golang.org/x/* packages
examples/main.go:19:58:use EXACTLY 2048 bits when generating RSA keys for some reason
examples/main.go:19:29:don't use math.Rand as source of entropy
examples/main.go:21:17:we don't use fmt.Errorf for some silly reaosn
examples/main.go:22:24:we don't use fmt.Errorf for some silly reaosn
examples/main.go:22:35:don't use uppercase error message string in fmt.Errorf formatted errors
examples/main.go:34:2:don't use http.Handle/http.HandleFunc which uses the DefaultServeMux due to possible side-effects
examples/main.go:35:2:don't use http.Handle/http.HandleFunc which uses the DefaultServeMux due to possible side-effects
examples/main.go:40:26:don't use old/weak tls cipher suites in your tls.Config
examples/main.go:45:38:don't use net.Listen to listen on all IPv4 addresses
examples/main.go:45:38:don't use net.Listen to listen on a random ports
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

### Struct

Using the `"struct"` type you can declare rules for structs.

#### Available Options for Struct

| Option         | Description                                             | Required  |
| ---------------|:--------------------------------------------------------|----------:|
| `name`         | the name of the package struct to inspect               | false     |
| `field`        | the specific struct field to inspect                    | false     |

```json
{
    "type": "struct",
    "comment": "don't use tls 1.3 ciphers in a tls.Config",
    "name": "tls.Config",
    "field": "CipherSuites",
    "cannot_match": [
        "tls.TLS_AES_128_GCM_SHA256",
        "tls.TLS_AES_256_GCM_SHA384",
        "TLS_CHACHA20_POLY1305_SHA256"
    ]
}
```
