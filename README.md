# asdev

`asdev` is a tool to ease with the development of ASUSTOR apps.

Currently `asdev` is experimental and only supports uploading APKs to update existing apps (i.e. submitting updates to App Central). It has not been tested for the creation of new applications.

`asdev` uses Google Chrome (or Chromium) to log on to the ASUSTOR Developer Corner and perform the necessary actions to update apps.

## Features

* Submit an updated APK to App Central
    * Reads changelog and description from APK and makes sure the fields are up-to-date
    * Re-applies all current app categories

## Requirements

* [Go](https://golang.org/dl/)
* [Google Chrome](https://www.google.com/chrome/browser/desktop/index.html) or [Chromium](https://www.chromium.org/getting-involved/download-chromium)
    * Other browsers that support the Chrome Debugging Protocol might work as well

## Installation

```console
$ go get -u github.com/mafredri/asdev
```

## Usage

The `asdev update` command is used to deploy one or more `.apk`'s to the ASUSTOR Developer Corner. For authentication, the username and password must be provided via environment variables or command line flags.

```console
$ export ASDEV_USERNAME="my-user"
$ export ASDEV_PASSWORD="my-password"
$ asdev update ./path/to/my.apk ./path/to/my-other.apk

OR:

$ asdev --username my-user --password my-password update ./path/to/my.apk ./path/to/my-other.apk
```

By default Chrome is run in headless mode, if you wish to see what `asdev` is doing, you can disable headless mode with the `--no-headless` command line flag.

**NOTE:** If you're not using macOS or using Chromium or Chrome Canary, you should provide the path to the browsers executable via the `--browser` cli flag (the environment variable `ASDEV_BROWSER` can also be used).

For more information, see `--help`:

```console
$ asdev --help-long
usage: asdev [<flags>] <command> [<args> ...]

Flags:
  -h, --help               Show context-sensitive help (also try --help-long and
                           --help-man).
  -u, --username=USERNAME  Username (for login)
  -p, --password=PASSWORD  Password (for login)
      --browser="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
                           Path to Chrome or Chromium executable
      --no-headless        Disable (Chrome) headless mode
      --timeout=10m        Command timeout
  -v, --verbose            Verbose mode

Commands:
  help [<command>...]
    Show help.


  show categories
    Show all available categories


  list
    List apps and status


  update [<flags>] <APKs>...
    Update apps by uploading one or more APK(s)

    -c, --category=CATEGORY ...  (NOT IMPLEMENTED) Change categorie(s)
    -t, --tag=TAG ...            (NOT IMPLEMENTED) Change tag(s)
    -b, --beta                   (NOT IMPLEMENTED) Beta app
    -i, --icon=ICON              (NOT IMPLEMENTED) Change icon (256x256)

  create --category=CATEGORY --tag=TAG [<flags>] <APKs>...
    (NOT IMPLEMENTED) Submit a new application by uploading one or more APK(s)

    -c, --category=CATEGORY ...  Set categorie(s)
    -t, --tag=TAG ...            Set tag(s)
    -b, --beta                   Set app to beta status
    -i, --icon=ICON              Set icon (256x256)
```

## Preview

```console
$ asdev list
+-------------+-------------+--------+------------+----------------------+---------------------+----------------+
| PACKAGE     | NAME        | ARCH   | VERSION    | CATEGORIES           | LAST UPDATE         | STATUS         |
+-------------+-------------+--------+------------+----------------------+---------------------+----------------+
| deluge      | Deluge      | x86-64 | 1.3.15-r1  | Download             | 2017-08-11 07:18:26 | Ready for sale |
| deluge      | Deluge      | arm    | 1.3.15-r1  | Download             | 2017-08-11 07:18:06 | Ready for sale |
| deluge      | Deluge      | i386   | 1.3.15-r1  | Download             | 2017-08-11 07:18:15 | Ready for sale |
| go          | Go          | x86-64 | 1.8.3      | Utility, Framework & | 2017-08-06 16:10:11 | In Review      |
|             |             |        |            | Library              |                     |                |
| go          | Go          | arm    | 1.8.3      | Utility, Framework & | 2017-08-06 16:27:53 | In Review      |
|             |             |        |            | Library              |                     |                |
| go          | Go          | i386   | 1.8.3      | Utility, Framework & | 2017-08-06 16:14:23 | In Review      |
|             |             |        |            | Library              |                     |                |
| jackett     | Jackett     | any    | 0.7.411    | Download             | 2016-10-08 12:11:33 | Ready for sale |
| qbittorrent | qBittorrent | x86-64 | 3.3.14     | Download             | 2017-08-11 07:18:36 | Ready for sale |
| qbittorrent | qBittorrent | arm    | 3.3.14     | Download             | 2017-08-11 07:18:54 | Ready for sale |
| qbittorrent | qBittorrent | i386   | 3.3.14     | Download             | 2017-08-11 07:19:12 | Ready for sale |
| radarr      | Radarr      | x86-64 | 0.2.0.453  | Download, Multimedia | 2017-03-13 21:02:51 | Ready for sale |
| radarr      | Radarr      | i386   | 0.2.0.453  | Download, Multimedia | 2017-03-13 21:02:20 | Ready for sale |
| radarr      | Radarr      | arm    | 0.2.0.453  | Download, Multimedia | 2017-03-13 17:59:53 | Ready for sale |
| sonarr      | Sonarr      | arm    | 2.0.0.4645 | Download             | 2017-03-13 21:11:54 | Ready for sale |
| sonarr      | Sonarr      | x86-64 | 2.0.0.4855 | Download, Multimedia | 2017-08-18 08:44:44 | Ready for sale |
| sonarr      | Sonarr      | i386   | 2.0.0.4855 | Download, Multimedia | 2017-08-18 08:44:23 | Ready for sale |
| tmux        | tmux        | arm    | 2.5.0      | Utility              | 2017-08-07 01:03:32 | In Review      |
| tmux        | tmux        | x86-64 | 2.5.0      | Utility              | 2017-08-07 01:03:59 | In Review      |
| tmux        | tmux        | i386   | 2.5.0      | Utility              | 2017-08-07 01:03:48 | In Review      |
+-------------+-------------+--------+------------+----------------------+---------------------+----------------+
```
