# mehr

Operating system-independent package and configuration manager

## Name

`mehr` is a pun on `nix`. `nix` in German means nothing, mehr translates to more, which is inherently more than nothing.

## Goals

The goal is to provide a singular interface for managing packages. This interfaces follows the following principles:

- high performance, bottle neck should be the underlying package managers
- minimal and well designed configuration file, see [TOML](https://toml.io/en/)
- support as many package managers and thus systems as possible
- automatically detect package managers
- support all package management interactions: install, upgrade, remove
- support configuring packages as well as installing them
- extensive documentation

## Supported Package Managers

| Package manager | Installing packages | Removing packages | Updating packages |
| --------------- | ------------------- | ----------------- | ----------------- |
| apt             | ✅                  | ✅                | ✅                |
| pacman          | ✅                  | ✅                | ✅                |
| cargo             | ✅                  | ✅                | ✅                |
| npm             | ✅                  | ✅                | ✅                |
| pnpm             | ✅                  | ✅                | ✅                |
| nix             | ❌                  | ❌                | ❌                |

## Documentation

### Command line interface

`mehr` works by accepting a sub comment and the sub comment accepting options
and flags.

#### Generate a default configuration

Due to the minimal nature of the `mehr` configuration the file can be hand
written. Still for the lazy, `mehr` supports generating a sensible default
configuration with extensive comments:

```shell
$ mehr init
```

#### Installing packages

> Installing packages via `mehr install <package>` does not modify the `mehr`
> configuration file, thus these packages are referred to as temporary
> packages. `mehr` does not manage dependencies of packages and performs
> install, upgrade and remove recursively, there always installing, updating and
> removing packages recursively.

Installing a singular package:

```shell
$ mehr install <package>
```

Installing multiple packages:

```shell
$ mehr install <package> <package>
```

Installing all packages defined in your `mehr` configuration file:

```shell
$ mehr install
```

#### Listing packages

Use `mehr list` to display a list of the currently installed packages:

```shell
$ mehr list
info: found 4 installed packages, 2 of them not specified in ~/.config/mehr/mehr.toml
kitty@0.21.2
firefox@119.0
info: temporary packages:
neovim@0.10.0-dev
falkon@3.2.0
```

Only list temporary packages:

```shell
$ mehr list --temporary
neovim@0.10.0-dev
falkon@3.2.0
$ mehr list -t
neovim@0.10.0-dev
falkon@3.2.0
```

Only list configured packages:

```shell
$ mehr list --permanent
kitty@0.21.2
firefox@119.0
$ mehr list -p
kitty@0.21.2
firefox@119.0
```

#### Upgrade packages

> Upgrading does not modify the `mehr` configuration, the latest version is
> used for upgrades.

To upgrade all packages on your system, run:

```shell
$ mehr upgrade
```

Updating a singular package:

```shell
$ mehr upgrade <package>
```

Updating multiple packages:

```shell
$ mehr upgrade <package> <package>
```

#### Removing packages

To remove all packages installed with `mehr` on your system, run:

```shell
$ mehr remove
```

Removing a singular package:

```shell
$ mehr remove <package>
```

Removing multiple packages:

```shell
$ mehr remove <package> <package>
```

#### Saving system state

Packages installed via `mehr install` are not reflected in the
`~/.config/mehr/mehr.toml`, but in `~/.config/mehr/lock.toml`. Saving
temporary installed packages into a configuration can be done via `mehr save`.
This reads the lock file and generates a new `mehr.toml` file into
`~/.config/mehr/` postfixed with the current time stamp.

```shell
$ mehr save
```

#### Restoring system state

`mehr sync` either forwards or resets the systems state to the state of the
configuration file, thus syncing both and removing the current lock file.

Packages not installed but found in the `mehr` configuration will be installed
upon running `mehr sync`. If packages are installed on the system but aren't
reflected in the configuration (referred to as temporary packages), the system
can be synced to the configuration via `mehr sync`:

```shell
$ mehr sync
info: detected 2 packages on your system that are not in ~/.config/mehr/mehr.toml:
neovim@0.10.0-dev
falkon@3.2.0
Are you sure your want to remove (2) temporary packages from your system? [Y/n]
```

Skip the confirmation prompt:

```shell
$ mehr sync --force
info: detected 2 packages on your system that are not in ~/.config/mehr/mehr.toml:
neovim@0.10.0-dev
falkon@3.2.0
info: removing 2 packages
```

#### Running commands

`mehr` accepts commands in its configuration file, these can be ran with `mehr run <cmd>`:

```shell
$ mehr run l
info: Got command: "ls -la" with cwd: "" and env: map[]
total 20
drwxr-xr-x 1 teo teo  118 Nov 16 12:15 .
drwxr-xr-x 1 teo teo 1006 Nov 13 13:04 ..
drwxr-xr-x 1 teo teo  126 Nov 16 12:11 cmd
drwxr-xr-x 1 teo teo   36 Nov 16 12:15 config
drwxr-xr-x 1 teo teo  204 Nov 16 12:13 .git
-rw-r--r-- 1 teo teo  242 Nov  6 09:04 go.mod
-rw-r--r-- 1 teo teo 1067 Nov  6 09:04 go.sum
drwxr-xr-x 1 teo teo   14 Nov  8 08:49 lock
drwxr-xr-x 1 teo teo   12 Nov  8 08:49 log
-rw-r--r-- 1 teo teo  109 Nov 16 10:36 main.go
drwxr-xr-x 1 teo teo   40 Nov 16 10:29 pkgmgr
-rw-r--r-- 1 teo teo 5573 Nov 16 12:15 README.md
drwxr-xr-x 1 teo teo   66 Nov 16 12:02 types
info: Ran command "l" successfully
```

View all available commands with `mehr run help`:

```shell
l:
        Command: "ls -la"

run:
        Command: "go run ."

generate:
        Command: "go generate"

build:linux/amd64:
        Command: "go build ."
        Env:
                GOOS: "linux"
                GOARCH: "amd64"
```

See Configuration/Commands for configuring commands.

### Configuration

`mehr` expects the `mehr.toml` file to be present at one of the following
locations:

| Operating system | Priority 1                        | Priority 2                                         |
| ---------------- | --------------------------------- | -------------------------------------------------- |
| Linux like       | `$XDG_CONFIG_HOME/mehr/mehr.toml` | `$HOME/mehr.toml`                                  |
| MacOS            | `$XDG_CONFIG_HOME/mehr/mehr.toml` | `$HOME/Library/Application Support/mehr/mehr.toml` |
| Windows          | `%AppData%/mehr/mehr.toml`        |                                                    |

#### Packages

#### System configuration

#### Commands
