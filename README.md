# `lovli-cli`

`lovli-cli` allows you to enjoy the [Lovli](https://lovli.fyi) URL shortening
service from the comfort of your terminal.

## Installation

> **Note on code signing for macOS users**
>
> The `lovli` binary is *unsigned*.
>
> To prevent Gatekeeper from blocking the execution of the program on macOS,
> you need to disable code signing enforcement on the `lovli` executable.
> The instructions below are meant to address this issue.
>
> Alternatively, you may consider building from source.

### Homebrew (macOS/Linux)

```sh
brew install --no-quarantine --cask callmesangio/lovli-cli/lovli-cli
```

### Manual download

- Download a release tarball from [here](https://github.com/callmesangio/lovli-cli/releases).
- Untar the archive.
- Run `xattr -dr "com.apple.quarantine" /path/to/lovli`.
- Copy the `lovli` binary to a directory in your `$PATH`.

## Usage

```sh
Usage of lovli:
  -h    Print help
  -s string
        Shorten the URL passed as argument
  -v    Print version
```

Example:

```sh
lovli -s https://example.com
```

## How to release

- Bump `version` string (`internal/app/app.go`).
- `git add . && git commit -m "Release X.Y.Z" && git push`
- `git tag X.Y.Z && git push --tags`
