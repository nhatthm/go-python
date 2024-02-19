# Go high-level bindings for the CPython-3 C-API

[![GitHub Releases](https://img.shields.io/github/v/release/nhatthm/go-python)](https://github.com/nhatthm/go-python/releases/latest)
[![Build Status](https://github.com/nhatthm/go-python/actions/workflows/test.yaml/badge.svg)](https://github.com/nhatthm/go-python/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/nhatthm/go-python/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/nhatthm/go-python)
[![Go Report Card](https://goreportcard.com/badge/go.nhat.io/python)](https://goreportcard.com/report/go.nhat.io/python)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/go.nhat.io/python)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

> [!IMPORTANT]
> **Currently supports python-3.11 only.**

The package provides a higher level for interacting with Python-3 C-API without directly using the `PyObject`.

Main goals:
- Provide generic a `Object` that interact with Golang types.
- Automatically marshal and unmarshal Golang types to Python types.

## Prerequisites

- `Go >= 1.22`
- `Python = 3.11`

## Install

```bash
go get go.nhat.io/python/v3
```

## Examples

```go
package main

import (
    "fmt"

    python3 "go.nhat.io/python/v3"
)

func main() {
    sys := python3.MustImportModule("sys")
    version := sys.GetAttr("version_info")

    pyMajor := version.GetAttr("major")
    defer pyMajor.DecRef()

    pyMinor := version.GetAttr("minor")
    defer pyMinor.DecRef()

    pyReleaseLevel := version.GetAttr("releaselevel")
    defer pyReleaseLevel.DecRef()

    major := python3.AsInt(pyMajor)
    minor := python3.AsInt(pyMinor)
    releaseLevel := python3.AsString(pyReleaseLevel)

    fmt.Println("major:", major)
    fmt.Println("minor:", minor)
    fmt.Println("release level:", releaseLevel)

    // Output:
    // major: 3
    // minor: 11
    // release level: final
}
```

```go
package main

import (
    "fmt"

    python3 "go.nhat.io/python/v3"
)

func main() { //nolint: govet
    math := python3.MustImportModule("math")

    pyResult := math.CallMethodArgs("sqrt", 4)
    defer pyResult.DecRef()

    result := python3.AsFloat64(pyResult)

    fmt.Printf("sqrt(4) = %.2f\n", result)

    // Output:
    // sqrt(4) = 2.00
}
```

## Donation

If this project help you reduce time to develop, you can give me a cup of coffee :)

### Paypal donation

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;or scan this

<img src="https://user-images.githubusercontent.com/1154587/113494222-ad8cb200-94e6-11eb-9ef3-eb883ada222a.png" width="147px" />
