# `factom-cli` Completion

Adds command line completion for
[`factom-cli`](https://github.com/FactomProject/factom-cli) for Bash, Zsh, and
Fish.

This program uses
[`github.com/posener/complete`](https://github.com/posener/complete) to
implement the completion.

## Build dependencies
- Golang 1.10 or higher

## Installation
```
go get -u github.com/AdamSLevy/complete-factom-cli
go install github.com/AdamSLevy/complete-factom-cli
complete-factom-cli -install -y
source ~/.bashrc
```

## Updating
```
go get -u github.com/AdamSLevy/complete-factom-cli
go install github.com/AdamSLevy/complete-factom-cli
```
You do not need to rerun the `-install` or `source` commands.
