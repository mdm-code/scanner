<h1 align="center">
  <div>
    <img
      src="https://raw.githubusercontent.com/mdm-code/mdm-code.github.io/main/scanner_logo.jpeg"
      alt="logo"
      style="object-fit: contain"
      width="40%"
    />
  </div>
</h1>

<h4 align="center">Custom Go text token scanner implementation</h4>

<div align="center">
<p>
    <a href="https://github.com/mdm-code/scanner/actions?query=workflow%3ACI">
        <img alt="Build status" src="https://github.com/mdm-code/scanner/workflows/CI/badge.svg">
    </a>
    <a href="https://app.codecov.io/gh/mdm-code/scanner">
        <img alt="Code coverage" src="https://codecov.io/gh/mdm-code/scanner/branch/main/graphs/badge.svg?branch=main">
    </a>
    <a href="https://opensource.org/licenses/MIT" rel="nofollow">
        <img alt="MIT license" src="https://img.shields.io/github/license/mdm-code/scanner">
    </a>
    <a href="https://goreportcard.com/report/github.com/mdm-code/scanner">
        <img alt="Go report card" src="https://goreportcard.com/badge/github.com/mdm-code/scanner">
    </a>
    <a href="https://pkg.go.dev/github.com/mdm-code/scanner">
        <img alt="Go package docs" src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white">
    </a>
</p>
</div>

Package `scanner` is a custom text scanner implementation. It has the same
idiomatic Go scanner programming interface, and it lets the client to freely
navigate the buffer. The scanner is also capable of peeking ahead of the
cursor. Read runes are rendered as tokens with additional information on their
position in the buffer. Consult the [package documentation](https://pkg.go.dev/github.com/mdm-code/scanner) or see
[Usage](#usage) to see how to use it.


## Installation

Use the following command to add the package to an existing project.

```sh
go get github.com/mdm-code/scanner
```


## Usage

Here is a snippet showing the basic usage of the scanner to read text as a stream
of tokens using the public API of the `scanner` package.

```go
package main

import (
    "bufio"
    "fmt"
    "log"
    "os"

    "github.com/mdm-code/scanner"
)

func main() {
    r := bufio.NewReader(os.Stdin)
    s, err := scanner.New(r)
    if err != nil {
        log.Fatalln(err)
    }
    var ts []scanner.Token
    for s.Scan() {
        t := s.Token()
        ts = append(ts, t)
    }
    fmt.Println(ts)
}
```


## Development

Consult [Makefile](Makefile) to see how to format, examine code with `go vet`,
run unit test, run code linter with `golint` in order to get test coverage and
check if the package builds all right.

Remember to install `golint` before you try to run tests and test the build:

```sh
go install golang.org/x/lint/golint@latest
```


## License

Copyright (c) 2023 Michał Adamczyk.

This project is licensed under the [MIT license](https://opensource.org/licenses/MIT).
See [LICENSE](LICENSE) for more details.
