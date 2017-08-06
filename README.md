# asdev

`asdev` is a tool to ease with the development of ASUSTOR apps.

## Status

* Experimental
* Uploading APKs to App Central

## Requirements

* Google Chrome or Chromium

## Installation

This tool is written in Go and requires Go to be installed (for building).

```console
$ go get -u github.com/mafredri/asdev
```

## Usage

Login via environment variables:

```console
$ export ASDEV_USERNAME="my-user"
$ export ASDEV_PASSWORD="my-password"
$ asdev -apk ./path/to/my.apk -apk ./path/to/my-other.apk
```

Or login via command line:

```console
$ asdev -username my-user -password my-password -apk ./path/to/my.apk -apk ./path/to/my-other.apk
```

**NOTE:** If the Chrome binary on your system exists in another location than the one expected by `asdev` (e.g. `/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`), then please provide the path to the browser via the command line option `-browser`.

```console
$ asdev -browser /my/path/to/chrome -apk ...
```