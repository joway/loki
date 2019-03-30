# Loki

![GitHub release](https://img.shields.io/github/tag/joway/loki.svg?label=release)
[![Go Report Card](https://goreportcard.com/badge/github.com/joway/loki)](https://goreportcard.com/report/github.com/joway/loki)
[![codecov](https://codecov.io/gh/joway/loki/branch/master/graph/badge.svg)](https://codecov.io/gh/joway/loki)
[![CircleCI](https://circleci.com/gh/joway/loki.svg?style=shield)](https://circleci.com/gh/joway/loki)

Inspired by a popular NodeJS logger [debug](https://www.npmjs.com/package/debug).

## Feature

- Colorful console output
- Enable logger by environment easily
- Minimalist

![](./demo.png)

## Install

```bash
go get github.com/joway/loki@latest
```

## API

[![Go Doc](https://godoc.org/github.com/joway/loki?status.svg)](https://godoc.org/github.com/joway/loki)

## Usage

### Create you own logger

```go
logger := loki.New("app:xxx")
logger.Info("x: %s", "hi")
logger.Debug("x: %s", "hi")
logger.Error("x: %s", "hi")
```

To enable the logger, just add env `LOKI_ENV=app:xxx` in your command, like: `LOKI_ENV=app:xxx ./main`.

You can also use `LOKI_ENV=app:a,app:b,model:a` to enable multiple logger.

`LOKI_ENV=*` can enable all loggers.

### Use root logger simply

It will not been affected by `LOKI_ENV`.

```go
loki.SetLevel(loki.INFO)

loki.Info("x: %s", "hi")
loki.Debug("x: %s", "hi")
loki.Error("x: %s", "hi")
```

### Use file handler

```go
// "02 Jan 06 15:04 MST"
l := loki.New("app:xxx")
fp, err := os.OpenFile("test.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
defer fp.Close()
flushIntervalMs := 1000
l.SetHandler(loki.NewFileHandler(fp, flushIntervalMs))

loki.Info("x: %s", "hi")
```

### Change time format

```go
// "02 Jan 06 15:04 MST"
loki.SetTimeFormatter(time.RFC822)

// disable time output
loki.SetTimeFormatter("")
```

### Use you own logger formatter

```go
type ErrFormatter struct {
	loki.Formatter
}

func (f ErrFormatter) format(a ...interface{}) string {
	err := a[0].(error)
	return fmt.Sprintf("Error %v", err)
}

logger := loki.New("app:xxx")
f := ErrFormatter{}
logger.SetFormatter(f)
logger.Debug(errors.New("test error"))
```
