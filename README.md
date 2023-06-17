
## srm (Safe ReMove)

<p>
  <img src="./assets/logo.png" width="128"/>
  <br>
</p>

`tl;dr:` ```srm``` is a command-line Safe ReMove files, that allows recover it if necessarily.

Deleted files are sent to (```~/.srm/```) giving you a chance to recover it. 

### Installation
   
Or get a binary [https://github.com/waldirborbajr/srm/releases][release] (Linux x86_64, ARMv7 and macOS), untar it, and move it somewhere on your $PATH:
```sh
#+BEGIN_EXAMPLE
$ tar xvzf srm-*.tar.gz
$ mv srm /usr/local/bin
#+END_EXAMPLE
```

### Usage

```sh
# Safe removing a file
$ srm srm file.txt

# Removing a file without safe restore enable
$ srm frm file.txt

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
- [ ] Safe Remove directory with content and subdirectories
- [ ] Restore directory with files and subdirectories
- [ ] Add short command option
- [ ] Implement tests
- [ ] TTL to purge files removed with safe option.
- [ ] ??? to be continued

## Technology

<img src="assets/gopher.png" alt="srm" width="32" /> 

[GO](https://go.dev/)
