# glint

> ⚠️ Under development, contributions are welcome!

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

| Option        | Description                                                              | Required  |
| ------------- |:-------------------------------------------------------------------------|----------:|
| `cannot_match`| array of regular expressions for packages the program **should not** use | false     |
| `must_match`  | array of regular expressions the packages program **should** use         | false     |

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
| `call_match`     | array of regular expressions to match against a package.Method call to inspect | false     |
| `argument`       | the index of the method argument (starting at `0` for the first) to inspect    | false     |
| `less_than`      | the method argument, if it's an integer, must be less than the given value     | false     |
| `greater_than`   | the method argument, if it's an integer, must be greater than the given value  | false     |
| `equals`         | the method argument, if it's an interger, must equal exactly the given value   | false     |
| `dont_use`       | the method call should not be used if the given faluse is true                 | false     |
| `cannot_match`   | array of regular expressions that method arguments *should not** match         | false     |

```json
{
   "type": "method",
   "comment": "use EXACTLY 2048 bits when generating RSA keys for some reason",
   "call": "rsa.GenerateKey",
   "argument": 1,
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

```json
{
    "type": "method",
    "comment": "potential for file inclusion in file name variable, be sure to clean the path if user input",
    "call_match": [
        "os.Open", "ioutil.ReadFile", "filepath.Join", "path.Join"
    ],
    "dont_use": true
}
```

```json
{
    "type": "method",
    "comment": "overly extensive file permissions detected",
    "call_match": [
        "os.Mkdir", "os.MkdirAll"
    ],
    "argument": 1,
    "less_than": 777
}
```

### Struct

Using the `"struct"` type you can declare rules for structs.

#### Available Options for Struct

| Option         | Description                                             | Required  |
| ---------------|:--------------------------------------------------------|----------:|
| `name`         | the name of the package struct to inspect               | false     |
| `field`        | the specific struct field to inspect                    | false     |

> **Note**: For the `field` option, if there are no fields defined when creating the struct in the inspected source code, then the assumed value to check against is `nil` for any type. This isn't the "zero value" you might expect, but greatly simplifies the config for checking structs.

```json
{
    "type": "struct",
    "comment": "don't use TLS 1.3 ciphers in a tls.Config",
    "name": "tls.Config",
    "field": "CipherSuites",
    "cannot_match": [
        "tls.TLS_AES_128_GCM_SHA256",
        "tls.TLS_AES_256_GCM_SHA384",
        "tls.TLS_CHACHA20_POLY1305_SHA256"
    ]
}
```

```json
{
    "type": "struct",
    "comment": "always set ReadTimeout when creating an http.Server",
    "name": "http.Server",
    "field": "ReadTimeout",
    "cannot_match": [
        "0", "nil"
    ]
}
```

```json
{
    "type": "struct",
    "comment": "http.Server ReadTimeout field should never be a single second",
    "name": "http.Server",
    "field": "ReadTimeout",
    "cannot_match": [
        "time.Second"
    ]
}
```

```json
{
    "type": "struct",
    "comment": "always set ReadTimeout when creating an ExampleServer",
    "name": "ExampleServer",
    "field": "ReadTimeout",
    "cannot_match": [
        "0", "nil"
    ]
}
```

> **Note**: The `ExampleServer` struct example is meant to show how to inspect non-imported packaged structs.

```go
type ExampleServer struct {
    ReadTimeout time.Duration
    // ...
}

// will catch this
var example = ExampleServer{
    ReadTimeout: 0,
}

// and this
var example2 = ExampleServer{}
```

```json
{
    "type": "struct",
    "comment": "for most structs, ReadTimeout should almost never be 0",
    "field": "ReadTimeout",
    "cannot_match": [
        "0"
    ]
}
```

> **Note**: The example above is meant to demonstrate that struct fields of any go struct type can be checked if the `name` option is ommitted. However, if the struct is created using the implicit zero value (field is not used during initialization), then this check will not apply to it unless it cannot match `nil` as well which could accidently be applied to any struct that has ignored fields and should generally be avoided.

```go
// it would be able to check this
ex1 := http.Server{
    ReadTimeout: 0,
}

// and this
ex2 := ExampleServer{
    ReadTimeout: 0,
}

// but not this
ex2 := ExampleServer{}

// or this this
ex2 := http.Server{}
```
