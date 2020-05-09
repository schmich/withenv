# withenv

`withenv` runs a command with environment variables merged from multiple sources.

## Install

[Download the zero-install binary](https://github.com/schmich/withenv/releases) to a directory on your `PATH`.

## Usage

```
Usage: withenv [--clear] --env=<file|@|->... [COMMAND [-- ARGS...]]

Run a command with environment from files - https://github.com/schmich/withenv

Arguments:
  COMMAND         Command to run.
  ARGS            Arguments to pass to command.

Options:
  -v, --version   Show the version and exit
  -c, --clear     Do not inherit current environment.
  -e, --env       Merge environment variables:
                  Use -e <file> to merge NAME=value entries from a file.
                  Use -e- to merge NAME=value entries from stdin.
                  Use -e @ to merge current environment.
```

## Environment Files

Environment files define one environment variable per line. Blank lines and lines starting with `#` (comments) are ignored. No shell execution, substitution, or interpolation is performed. See the [notes](#notes) section for details.

### Example

```bash
# Comments are ignored.
# Blank lines are also ignored.
# Values are handled verbatim as literal strings.

# Simple environment variable.
HELLO=world

# Empty value.
EMPTY_VAR=

# Spaces are allowed.
WITH_SPACES=A value with spaces, very cool.

# These are still ultimately strings.
SOME_INT=123
SOME_BOOL=false
SOME_NULL=null

# No environment substitution occurs.
NOT_SUBSTITUTED=$USER %USER%

# Shell commands aren't run.
JUST_A_STRING=$(pwd)
ALSO_A_STRING=`/bin/ls -lh`

# Single and double quotes are not removed.
QUOTED="Quotes are not removed from values"
SINGLE='Single quotes are not removed either'
ESCAPED=This \"is\" (''also'') \\`verbatim`\\

# AWS credentials.
AWS_ACCESS_KEY_ID=4cc355k3yid
AWS_SECRET_KEY=s0m3s3cr37
```

## Examples

The following examples assume the following:

Current environment:
```
SOURCE=current
SHELL=/bin/bash
```

Environment file `.env1`:
```
SOURCE=env1
COLOR=blue
```

Environment file `.env2`:
```
SOURCE=env2
VOLUME=max
```

**Example:** Show merged environment after loading an environment file (useful for debugging):

```bash
$ withenv -e .env1
SHELL=/bin/bash
SOURCE=env1
COLOR=blue
```

**Example:** Merge multiple environments from multiple files:

```bash
$ withenv -e .env1 -e .env2
SHELL=/bin/bash
SOURCE=env2
COLOR=blue
VOLUME=max
```

**Example:** Do not inherit current environment. Use environment exclusively from file:

```bash
$ withenv -c -e .env1
SOURCE=env1
COLOR=blue
```

**Example:** Use environment from stdin:

```bash
$ cat .env1 | withenv -e-
SHELL=/bin/bash
SOURCE=env1
COLOR=blue
```

**Example:** Use environment from file with current environment taking precedence:

```bash
$ withenv -c -e .env1 -e @
SHELL=/bin/bash
SOURCE=current
COLOR=blue
```

**Example:** Run a complex command with a merged environment:

```bash
$ withenv --clear --env .env.shared --env .env.production php -c /etc/php -f MyScript.php > out.log
```

## Notes

- Environment variables are case-sensitive
- Environment variables match the pattern `[A-Za-z_][A-Za-z0-9_]*`
- The target process is started using [`syscall.Exec`](https://golang.org/pkg/syscall/#Exec), see [execve(2)](https://linux.die.net/man/2/execve)

## License

Copyright &copy; 2020 Chris Schmich  \
MIT License. See [LICENSE](LICENSE) for details.
