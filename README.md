
## srm

<p>
  <img src="./assets/logo.png" width="128"/>
  <br>
</p>

rtl;dr:` srm it is a CLI (Command Line Interface), thats aims to be a Safe ReMove files, that allows recover it if necessarily.

### How to use

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
