package runnlib

import (
	"io"
	"os"
)

type CommandConfig struct {
	Environment map[string]string `json:"environment"`
	WorkDir     string            `json:"workdir"`
	User        string            `json:"user"`
	Cmd         string            `json:"cmd"`
	Args        []string          `json:"args"`
	Interpreter []string          `json:"interpreter"`
	Script      bool              `json:"script"`
	Interactive bool              `json:"interactive"`
	Stdout      string            `json:"stout"`
	Stderr      string            `json:"stderr"`
}

func GetOutput(config CommandConfig) (stdout io.Writer, stderr io.Writer, err error) {

	if config.Stderr == "stderr" {
		stderr = os.Stderr
	} else if config.Stderr != "" {
		stderr, err = os.Create(config.Stderr)
	}

	if config.Stdout == "stdout" {
		stdout = os.Stdout
	} else if config.Stdout != "" {
		stdout, err = os.Create(config.Stdout)
	}

	return
}
