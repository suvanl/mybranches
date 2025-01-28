# mybranches

![](https://github.com/user-attachments/assets/aee520f2-d710-4014-babb-09aa34fb2d5d)

A super simple program that allows you to interactively switch to a local branch matching a certain pattern. By default, this pattern is your username.

## Why does this exist?

While working on projects with multiple collaborators, I'd often find myself repeating these steps when trying to find one of my branches:

- have a number of local branches
- run `git mybranches` (alias for `git branch --list "<username>*"`)
- manually select + copy the name of the branch I'm looking for to the clipboard
- run `git switch <branchname>`

This automates that process.

## Installation
### Build from source
1. Clone this repository
2. Run `make -B`
3. Create symlink to directory on `$PATH` (use `echo $PATH` to check). 
    - E.g., if /usr/local/bin is on your PATH:
    - `ln -s ~/path/to/repo/bin/mybranches /usr/local/bin/mybranches`

## Usage
```
mybranches
```

**Optional flags**:
- `--pattern`: specify a custom pattern. This gets passed to `git branch --list <pattern>*`. Defaults to your system username.

> [!IMPORTANT]
> The "copy to clipboard" feature is currently only supported on macOS.
