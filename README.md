# mybranches

A super simple program that allows you to interactively switch to a local branch matching a certain pattern. By default, this pattern is your username.

##  Why does this exist?

I'd often find myself repeating these steps:

- have a number of local branches
- run `git mybranches` (alias for `git branch --list "<username>*"`)
- manually select + copy the name of the branch I'm looking for to the clipboard
- run `git switch <branchname>`

This automates that process.
