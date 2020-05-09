package main

import (
	"bufio"
	"fmt"
	logging "log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"github.com/jawher/mow.cli"
	"github.com/pkg/errors"
)

var version string
var commit string
var log *logging.Logger

func init() {
	log = logging.New(os.Stderr, "", 0)
}

func loadEnv(filename string) (map[string]string, error) {
	var err error
	var file *os.File

	if filename == "-" {
		file = os.Stdin
	} else {
		filename, err = filepath.Abs(filename)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to get absolute path for %s", filename))
		}

		file, err = os.Open(filename)
		if err != nil {
			return nil, errors.Wrap(err, "failed to open environment file")
		}

		defer file.Close()
	}

	entry := regexp.MustCompile(`(?i)^([a-z_][a-z0-9_]*)=(.*)$`)
	ignored := regexp.MustCompile(`(^\s*$)|(^#.*$)`)

	env := make(map[string]string)

	lineNum := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineNum++

		line := scanner.Text()
		if ignored.MatchString(line) {
			continue
		}

		match := entry.FindStringSubmatch(line)
		if len(match) == 0 {
			return nil, fmt.Errorf("environment file syntax error %s:%d", filename, lineNum)
		}

		name := match[1]
		value := match[2]

		env[name] = value
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return env, nil
}

func toMap(entries []string) map[string]string {
	env := make(map[string]string)

	for _, entry := range entries {
		pair := strings.SplitN(entry, "=", 2)
		name := pair[0]
		value := pair[1]
		env[name] = value
	}

	return env
}

func toEntries(env map[string]string) []string {
	entries := []string{}

	for name, value := range env {
		entry := name + "=" + value
		entries = append(entries, entry)
	}

	return entries
}

func mergeEnv(current map[string]string, overlay map[string]string) map[string]string {
	env := make(map[string]string)

	for name, value := range current {
		env[name] = value
	}

	for name, value := range overlay {
		env[name] = value
	}

	return env
}

func main() {
	app := cli.App("withenv", "Run a command with environment from file - https://github.com/schmich/withenv")

	app.Spec = "-e=<file|@|->... [COMMAND [-- ARGS...]]"

	app.Version("v version", "withenv "+version+" "+commit)

	sources := app.StringsOpt("e env", []string{}, "Environment file containing NAME=VALUE entries.\nUse -e- to read from stdin.\nUse -e @ to read current environment.")
	cmd := app.StringArg("COMMAND", "", "Command to run.")
	args := app.StringsArg("ARGS", []string{}, "Arguments to pass to command.")

	app.Action = func() {
		var err error

		current := toMap(os.Environ())
		combined := mergeEnv(map[string]string{}, current)

		for _, source := range *sources {
			var overlay map[string]string
			if source == "@" {
				overlay = current
			} else {
				if overlay, err = loadEnv(source); err != nil {
					log.Fatal(err)
				}
			}

			combined = mergeEnv(combined, overlay)
		}

		env := toEntries(combined)

		if *cmd == "" {
			for _, entry := range env {
				fmt.Println(entry)
			}
			return
		}

		command, err := exec.LookPath(*cmd)
		if err != nil {
			log.Fatal(err)
		}

		arguments := append([]string{command}, *args...)

		if err = syscall.Exec(command, arguments, env); err != nil {
			log.Fatal(err)
		}
	}

	app.Run(os.Args)
}
