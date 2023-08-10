## srm (Safe ReMove)

[![Go](https://github.com/waldirborbajr/srm/actions/workflows/go.yml/badge.svg)](https://github.com/waldirborbajr/srm/actions/workflows/go.yml) .::. [![golangci-lint](https://github.com/waldirborbajr/srm/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/waldirborbajr/srm/actions/workflows/golangci-lint.yml) .::. [![Typo Check](https://github.com/waldirborbajr/srm/actions/workflows/typo-check.yaml/badge.svg)](https://github.com/waldirborbajr/srm/actions/workflows/typo-check.yaml)

<p>
  <img src="./assets/logo.png" width="128"/>
  <br>
</p>

`tl;dr:` ```srm``` is a command-line Safe ReMove file, that allows recovery if necessarily.

Deleted files are sent to (```~/.srm/```) giving you a chance to recover it. 

### Installation
   
Or get a binary [https://github.com/waldirborbajr/srm/releases][release] (Linux x86_64, ARMv7, and macOS), untar it, and move it somewhere on your $PATH:

```sh
$ tar xvzf srm-*.tar.gz
$ mv srm /usr/local/bin
```

### Overview

```sh
srm
Usage: srm command [options]

srm - Safe ReMove is a simple tool to remove file/directory safety.

Option:

Commands:
srm - Remove a file/directory using the safe mode that preserves the file that is possible to restore
rst - Restore a file/directory that was deleted with a safe option
cls - Cleanup removed files after 18 days if not informed another day as parameter
hlp - Display this help information
ver - Prints version info to the console
```

```sh
srm srm
Usage: srm srm [options....]

Usage:
        srm srm --save file.bak
        srm srm --force file.bak
        srm srm --save file*

Options:
        -s, --safe Save removed file/directory for restore
        -f, --force Remove file/directory without restore option
```      

### Usage

```sh
# Safe removing a single file
$ srm srm -s file.txt

# Removing a single file without safe option
$ srm srm -f file.txt

# Safe removing with wildcard pattern
$ srm srm -s file*

# Restoring a file removed with safe option
$ srm rst file.txt

# Print version of srm
$ srm ver 
```

## License

[MIT](https://github.com/waldirborbajr/srm/blob/main/LICENSE)

## Legal

Copyright 2023 Waldir Borba Junior (<mailto:wborbajr@gmail.com>)
SPDX-License-Identifier: Apache-2.0

## TODO
- [x] On Safe Remove must compress "deleted" content to optimize storaged space
- [x] Safe Remove directory with content and subdirectories
- [ ] Restore directory with files and subdirectories
- [ ] Add short command option
- [ ] Implement tests
- [x] TTL to purge files removed with safe option.
- [x] Delete passing with wildcard pattern
- [ ] Open to PR and **new** features

## Technology

<img src="assets/gopher.png" alt="srm" width="32" /> 

[GO](https://go.dev/)
