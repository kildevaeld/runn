package runn

import (
	"errors"
	"io"
	"os"
	"os/exec"

	"github.com/kildevaeld/runn/runnlib"
)

const (
	Stdout = "stdout"
	Stderr = "stderr"
)

type Command struct {
	config runnlib.CommandConfig
}

func getCommand(config runnlib.CommandConfig) (cmd *exec.Cmd, err error) {

	var exe string
	var args []string

	if config.Cmd == "" {
		return nil, errors.New("no command")
	}

	if config.Script {
		if len(config.Interpreter) == 0 {
			return nil, errors.New("cannot run as script, when interpreter is set")
		}

		//args = []string{config.Cmd}
	}

	if len(config.Interpreter) > 0 {
		exe, err = exec.LookPath(config.Interpreter[0])
		if err != nil {
			return nil, err
		}
		if len(config.Interpreter) > 1 {
			args = config.Interpreter[1:]
		}

		args = append(args, config.Cmd)
		args = append(args, config.Args...)
	} else {
		exe, err = exec.LookPath(config.Cmd)
		if err != nil {
			return nil, err
		}

		args = config.Args
	}

	if exe == "" {
		return nil, errors.New("could not find exe")
	}

	cmd = exec.Command(exe, args...)
	cmd.Dir = config.WorkDir

	if cmd.Dir == "" {
		cmd.Dir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	env := os.Environ()
	if config.Environment != nil {
		for k, v := range config.Environment {
			env = append(env, k+"="+v)
		}
	}
	cmd.Env = env

	if config.Interactive {
		cmd.Stdin = os.Stdin
	}

	return cmd, nil
}

func (self *Command) Run() error {

	stdout, stderr, err := runnlib.GetOutput(self.config)

	if err != nil {
		return err
	}

	close := func(w io.Writer) {
		if w != os.Stdout && w != os.Stderr {
			if file, ok := w.(*os.File); ok {
				file.Close()
			}
		}
	}

	defer close(stdout)
	defer close(stderr)

	cmd, cerr := getCommand(self.config)
	cmd.Stderr = stderr
	cmd.Stdout = stdout

	if cerr != nil {
		return cerr
	}

	return cmd.Run()
}

func Cmd(config runnlib.CommandConfig) *Command {
	return &Command{config}
}
