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
	filename, err := filepath.Abs(filename)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to get absolute path for %s", filename))
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open environment file")
	}

	defer file.Close()

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

func sliceToMap(environ []string) map[string]string {
	env := make(map[string]string)

	for _, entry := range environ {
		pair := strings.SplitN(entry, "=", 2)
		name := pair[0]
		value := pair[1]
		env[name] = value
	}

	return env
}

func mapToSlice(env map[string]string) []string {
	environ := []string{}

	for name, value := range env {
		entry := name + "=" + value
		environ = append(environ, entry)
	}

	return environ
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
	app := cli.App("withenv", "Run command with environment from file.")

	app.Spec = "ENVFILE [COMMAND [-- ARGS...]]"

	app.Version("v version", "withenv "+version+" "+commit)

	envFile := app.StringArg("ENVFILE", "", "Environment file containing NAME=VALUE entries.")
	cmd := app.StringArg("COMMAND", "", "Command to run.")
	_ = app.StringsArg("ARGS", nil, "Arguments to pass to command.")

	app.Action = func() {
		current := sliceToMap(os.Environ())

		overlay, err := loadEnv(*envFile)
		if err != nil {
			log.Fatal(err)
		}

		env := mapToSlice(mergeEnv(current, overlay))

		if *cmd == "" {
			for _, entry := range env {
				fmt.Println(entry)
			}
		} else {
			command, err := exec.LookPath(*cmd)
			if err != nil {
				log.Fatal(err)
			}

			args := os.Args[2:]
			if err := syscall.Exec(command, args, env); err != nil {
				log.Fatal(err)
			}
		}
	}

	app.Run(os.Args)
}
