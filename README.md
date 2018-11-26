# pdump
## Description

Detect DDOS and execute arbitrary command

## Usage

```
Usage of pdump:
  -a uint
        alert threshould(Short) (default 10)
  -alert uint
        alert threshould (default 10)
  -b uint
        BufflerLength(Short)
  -buffer uint
        BufflerLength
  -e string
        exec command(Short)
  -exec string
        exec command
  -i uint
        monitor interval(Short) (default 30)
  -interval uint
        monitor interval (default 30)
  -n string
        monitor nic(Short)
  -nic string
        monitor nic
  -s uint
        monitor sec(Short) (default 5)
  -sec uint
        monitor sec (default 5)
  -version
        Print version information and quit.
```

## Install

To install, use `go get`:

```bash
$ go get -d github.com/pyama86/pdump
```

## Contribution

1. Fork ([https://github.com/pyama86/pdump/fork](https://github.com/pyama86/pdump/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[pyama86](https://github.com/pyama86)
