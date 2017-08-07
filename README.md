# asdev

`asdev` is a tool to ease with the development of ASUSTOR apps.

Currently `asdev` is experimental and only supports uploading APKs to update existing apps (i.e. submitting updates to App Central). It has not been tested for the creation of new applications.

`asdev` uses Google Chrome (or Chromium) to log on to the ASUSTOR Developer Corner and perform the necessary actions to update apps.

## Features

* Submit an updated APK to App Central
    * Reads changelog and description from APK and makes sure the fields are up-to-date
    * Re-applies all current app categories

## Requirements

* Google Chrome or Chromium

## Installation

This tool is written in Go and requires Go to be installed (for building).

```console
$ go get -u github.com/mafredri/asdev
```

## Usage

To upload files, a username and password must be provided (for authentication to the ASUSTOR Developer Corner).

Login (and upload) via environment variables:

```console
$ export ASDEV_USERNAME="my-user"
$ export ASDEV_PASSWORD="my-password"
$ asdev upload ./path/to/my.apk ./path/to/my-other.apk
```

Or provide login info via command line flags:

```console
$ asdev --username my-user --password my-password upload ./path/to/my.apk ./path/to/my-other.apk
```

By default Chrome is run in headless mode, if you wish to see what `asdev` is doing, you can disable headless mode with the `--no-headless` command line flag.

**NOTE:** If the Chrome binary on your system exists in another location than the one expected by `asdev` (e.g. `/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`), then please provide the path to the browser via the command line option `-browser`.

For more information, see `--help`:

```console
$ asdev --help
usage: asdev [<flags>] <command> [<args> ...]

Flags:
  -h, --help               Show context-sensitive help (also try --help-long and
                           --help-man).
  -u, --username=USERNAME  Username (for login)
  -p, --password=PASSWORD  Password (for login)
      --browser="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
                           Path to Chrome or Chromium binary
      --no-headless        Disable (Chrome) headless mode
  -v, --verbose            Verbose mode

Commands:
  help [<command>...]
    Show help.

  upload [<flags>] <APKs>...
    Upload one or multiple APK(s)
```
