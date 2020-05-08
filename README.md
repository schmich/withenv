# withenv

`withenv` runs a command with an environment comprising current environment variables and environment variables loaded from an environment file containing `NAME=value` entries.

## Install

[Download the zero-install binary](https://github.com/schmich/withenv/releases) to a directory on your `PATH`.

## Usage

```
Usage: withenv ENVFILE [COMMAND [-- ARGS...]]

Run a command with environment from file - https://github.com/schmich/withenv

Arguments:
  ENVFILE         Environment file containing NAME=VALUE entries.
  COMMAND         Command to run.
  ARGS            Arguments to pass to command.

Options:
  -v, --version   Show the version and exit
```

## Examples

Show merged environment after loading environment file without running any command:

```
withenv .env
```

Run command with loaded environment:

```
withenv .env php MyScript.php
```

## Notes

- Environment variables are case-sensitive
- The process is started using [`syscall.Exec`](https://golang.org/pkg/syscall/#Exec)

## License

Copyright &copy; 2020 Chris Schmich  \
MIT License. See [LICENSE](LICENSE) for details.
