package runn

import (
	"bytes"
	"errors"
	"net"
	"strings"
	"text/template"

	"github.com/fatih/camelcase"
	"github.com/kildevaeld/dict"
	"github.com/kildevaeld/runn/runnlib"
)

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {

			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

type Bundle struct {
	workdir string
	config  runnlib.Bundle
}

func camelCaseToEnv(k string) string {
	split := camelcase.Split(k)
	for i, s := range split {
		split[i] = strings.ToUpper(s)
	}
	return strings.Join(split, "_")
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

	if cmd.Environment != nil {
		out.Environment = make(map[string]string)

		for k, v := range cmd.Environment {
			out.Environment[camelCaseToEnv(k)], _ = interpolateString(v, locals)
		}
		for k, v := range locals {
			if s, ok := v.(string); ok {
				out.Environment[camelCaseToEnv(k)] = s
			}

		}
	}

	return out, nil
}

func getCommandInBundle(bundle runnlib.Bundle, name string) *runnlib.BundleCommand {
	for _, c := range bundle.Commands {
		if c.Name == name {
			return &c //&c.Command
		}
	}
	return nil
}

func (self *Bundle) Run(name string, config RunConfig) error {

	comm := getCommandInBundle(self.config, name)
	if comm == nil {
		//fmt.Printf("%#v\n", self.config)
		return errors.New("No cmd: " + name)
	}

	locals := dict.NewMap()
	locals["WorkDir"] = self.workdir
	locals["HostIP"] = GetLocalIP()
	locals["RuntimePath"] = config.RuntimePath
	if config.Locals != nil {
		for k, v := range config.Locals {
			locals[k] = v
		}
	}

	//config.Environ = append(Config.Environ, "RUNTIME_PATH="+config.RuntimePath)
	//config.Environ.Add("RUNTIME_PATH=" + config.RuntimePath)

	conf, err := interpolateCommand(&comm.Command, locals)
	if err != nil {
		return err
	}

	cmd := Cmd(self.config.Name, comm.Name, conf)

	return cmd.Run(config)

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
