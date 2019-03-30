# Loki

![GitHub release](https://img.shields.io/github/tag/joway/loki.svg?label=release)
[![Go Report Card](https://goreportcard.com/badge/github.com/joway/loki)](https://goreportcard.com/report/github.com/joway/loki)
[![codecov](https://codecov.io/gh/joway/loki/branch/master/graph/badge.svg)](https://codecov.io/gh/joway/loki)
[![CircleCI](https://circleci.com/gh/joway/loki.svg?style=shield)](https://circleci.com/gh/joway/loki)

## Install

```bash
go get github.com/joway/loki@latest
```

## Usage

```go
package main

import "github.com/joway/loki"

func main() {
	loki.SetLevel(loki.INFO)

	loki.Info("x: %s", "hi")
	loki.Debug("x: %s", "hi")
	loki.Error("x: %s", "hi")	
}
```
