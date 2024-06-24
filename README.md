# errapy

## Overview
errapy is a Go library designed to streamline error handling with detailed error reporting. This library allows for enhanced error descriptions, making debugging and logging more efficient and informative.

## Features
- **Detailed Error Reporting**: Provides comprehensive error details.
- **Customizable Error Messages**: Tailor error messages to fit specific needs.
- **Easy Integration**: Simple and intuitive to integrate with existing projects.

## Installation
To install Errapy, use `go get`:
```bash
go get github.com/d3v3us/errapy
```

## Applying policies
```bash
package main

import (
    "fmt"
    "github.com/d3v3us/errapy"
)

func main() {
    policy := errapy.Policy(
        WithClassesRequired(true),
        WithCodesRequired(true),
    )
    err := errapy.New("Something went wrong", errapy.WithPolicy(policy))
    fmt.Println(err)
}
```