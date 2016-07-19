package runn

import (
	"bytes"
	"errors"
	"text/template"

	"github.com/kildevaeld/dict"
	"github.com/kildevaeld/runn/runnlib"
)

type Bundle struct {
	workdir string
	config  runnlib.Bundle
}

func interpolateString(arg string, locals dict.Map) (string, error) {
	tmpl, err := template.New("test").Parse(arg)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(nil)

	err = tmpl.Execute(buf, locals)
	if err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil

}

func interpolateStringSlice(args []string, locals dict.Map) ([]string, error) {
	var out []string
	for _, r := range args {
		if str, err := interpolateString(r, locals); err != nil {
			return nil, err
		} else {
			out = append(out, str)
		}
	}
	return out, nil
}

func interpolateCommand(cmd *runnlib.CommandConfig, locals dict.Map) (runnlib.CommandConfig, error) {
	out := runnlib.CommandConfig{}
	var err error
	if out.Args, err = interpolateStringSlice(cmd.Args, locals); err != nil {
		return out, err
	}

	if out.Cmd, err = interpolateString(cmd.Cmd, locals); err != nil {
		return out, err
	}
	if out.Interpreter, err = interpolateStringSlice(cmd.Interpreter, locals); err != nil {
		return out, err
	}
	if out.WorkDir, err = interpolateString(cmd.WorkDir, locals); err != nil {
		return out, err
	}

	if out.Stderr, err = interpolateString(cmd.Stderr, locals); err != nil {
		return out, err
	}
	if out.Stdout, err = interpolateString(cmd.Stdout, locals); err != nil {
		return out, err
	}

	return out, nil
}

func getCommandInBundle(bundle runnlib.Bundle, name string) *runnlib.CommandConfig {
	for _, c := range bundle.Commands {
		if c.Name == name {
			return &c.Command
		}
	}
	return nil
}

func (self *Bundle) Run(name string) error {

	comm := getCommandInBundle(self.config, name)
	if comm == nil {
		//fmt.Printf("%#v\n", self.config)
		return errors.New("No cmd: " + name)
	}

	locals := dict.NewMap()
	locals["WorkDir"] = self.workdir

	config, err := interpolateCommand(comm, locals)
	if err != nil {
		return err
	}

	cmd := Cmd(config)

	return cmd.Run()

	//return nil
}

func NewBundle(path string) (bundle *Bundle, err error) {

	/*if bs, err = ioutil.ReadFile(filepath.Join(path, "bundle.yaml")); err != nil {
		return nil, err
	} else if bs, err = ioutil.ReadFile(filepath.Join(path))*/
	var b runnlib.Bundle
	if err := runnlib.GetBundleFromPath(path, &b); err != nil {
		return nil, err
	}

	bundle = &Bundle{path, b}

	return bundle, err
}