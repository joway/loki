# Loki

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
