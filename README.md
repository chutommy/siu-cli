# Siu

Siu is a CLI application which can be launched with only three keystrokes "siu". It allows you to use shortcuts to start a browser with a specific URL. If a browser window is already opened, it only creates a new tab.

## Installation

Requires `go v1.13` or higher installed (https://golang.org/dl/).

```bash
go get github.com/chutified/siu
```

## Usage

Run 'siu -h` for the help.

command|action
-------|------
`siu`|Opens the URL with motions (supports multiple inputs)
`siu list`|Lists all motions
`siu set del [motion]`|Deletes one or multiple motions
`siu set new [motion]`|Creates a new motion
`siu set upd [motion]`|Updates the motion

All __motions__ can be identified by its _name_, _url_, _shortcut_ or _id_.

## Screenshots of CLI

### `siu -h` command:

![screenshot of siu --help](https://raw.githubusercontent.com/chutified/siu/master/img/00_siu_help.png)

### `siu list` command:

![screenshot of siu list](https://raw.githubusercontent.com/chutified/siu/master/img/01_siu_list.png)

### `siu` command (launch new browser):

![screenshot of siu run](https://raw.githubusercontent.com/chutified/siu/master/img/02_siu_run.png)
![screenshot of opened tabs](https://raw.githubusercontent.com/chutified/siu/master/img/03_siu_browser.png)

### `siu -h` command (browser exists, open new tab):

![screenshot of siu run](https://raw.githubusercontent.com/chutified/siu/master/img/04_siu_run.png)
![screenshot of opened tabs](https://raw.githubusercontent.com/chutified/siu/master/img/05_siu_browser.png)
