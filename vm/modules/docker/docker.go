//go:generate go-bindata -pkg docker -o docker_impl.go docker.async.js builder.js
package docker

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/kildevaeld/notto"
	"github.com/kildevaeld/notto/loop"
	"github.com/robertkrimen/otto"
)

type simple_task struct {
	id     int64
	err    error
	result otto.Value
	call   otto.Value
}

func (self *simple_task) SetID(id int64) { self.id = id }
func (self *simple_task) GetID() int64   { return self.id }

func (self *simple_task) Execute(vm *otto.Otto, loop *loop.Loop) error {

	var arguments []interface{}

	if self.err != nil {
		e, err := vm.Call(`new Error`, nil, self.err.Error())
		if err != nil {
			return err
		}

		arguments = append(arguments, e)
	} else {
		arguments = append(arguments, otto.NullValue())
	}

	/*if self.result != nil {
		if v, e := vm.ToValue(t.content); e != nil {
			return e
		} else {
			arguments = append(arguments, v)
		}

	}*/
	arguments = append(arguments, self.result)

	//arguments = append([]interface{}{self.call}, arguments...)
	if _, err := self.call.Call(otto.NullValue(), arguments...); err != nil {
		return err
	}
	/*if _, err := vm.Call(`Function.call.call`, nil, arguments...); err != nil {
		return err
	}**/

	return nil
}

func (self *simple_task) Cancel() {

}

func SimpleTask(vm *notto.Notto, cb otto.Value, worker func() (interface{}, error)) {

	task := &simple_task{call: cb}
	vm.Runloop().Add(task)

	go func() {
		defer vm.Runloop().Ready(task)
		var (
			i interface{}
			e error
			//v otto.Value
		)
		i, e = worker()
		task.err = e
		if i == nil {
			task.result = otto.UndefinedValue()
		} else {
			task.result, e = vm.ToValue(i)
			if e != nil {
				task.err = e
			}
		}

	}()

}

func Define(vm *notto.Notto) error {

	o, e := vm.Object("({})")
	if e != nil {
		return e
	}
	//err := vm.Set("__private_docker", privateDockerCall(vm))
	o.Set("docker", func(call otto.FunctionCall) otto.Value {

		str := call.Argument(0).String()
		privateDocker(vm, str, call.Argument(1))

		return otto.UndefinedValue()
	})
	val, _ := o.Get("docker")

	m := make(map[string]otto.Value)
	m["__private_docker"] = val
	loader := notto.CreateLoaderFromSourceWithPrivate(string(MustAsset("docker.async.js")), "", m)

	vm.AddModule("docker.async", loader)
	loader = notto.CreateLoaderFromSource(string(MustAsset("builder.js")), "")
	vm.AddModule("docker.builder", loader)
	return nil
}

func Define2(vm *notto.Notto) error {

	o, e := vm.Object("({})")
	if e != nil {
		return e
	}
	//err := vm.Set("__private_docker", privateDockerCall(vm))
	o.Set("create", func(call otto.FunctionCall) otto.Value {
		var (
			e error
			c *docker_p
			v otto.Value
		)
		if c, e = createDocker(vm); e != nil {
			vm.Throw("DockerError", e)
		}
		if v, e = vm.ToValue(c); e != nil {
			vm.Throw("DockerError", e)
		}
		return v
	})

	vm.AddModule("docker", notto.CreateLoaderFromValue(o.Value()))

	return nil
}

func privateDocker(vm *notto.Notto, str string, call otto.Value) {
	SimpleTask(vm, call, func() (interface{}, error) {

		cmd := exec.Command("docker")
		stdout := bytes.NewBuffer(nil)
		stderr := bytes.NewBuffer(nil)
		conf := vm.ProcessAttr()

		cmd.Stderr = stderr
		cmd.Stdout = stdout
		cmd.Dir = conf.Cwd
		cmd.Args = strings.Split(str, " ")
		err := cmd.Run()

		if err != nil {
			return nil, fmt.Errorf("CMD: %s failed\n%s", str, stderr.Bytes())
		}

		if !cmd.ProcessState.Success() {
			return nil, errors.New(cmd.ProcessState.String())
		}

		out := strings.TrimSpace(string(stdout.Bytes()))

		if stderr.Len() > 0 {
			err = fmt.Errorf("CMD: %s failed\n%s", str, stderr.Bytes())
		}

		return out, err
	})
}

/*func privateDockerCall(vm *notto.Notto) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {

		var (
			cmd string
			err error
		)

		if len(call.ArgumentList) != 4 {
			vm.Throw("ParameterError", "Invalid number of arguments")
		}

		if !call.Argument(0).IsObject() {
			vm.Throw("ParameterError", "First argument must be an object")
		}

		if !call.Argument(1).IsString() {
			vm.Throw("ParameterError", "Second argument must be a string")
		} else {
			if cmd, err = call.Argument(1).ToString(); err != nil {
				vm.Throw("ParameterError", err)
			}
		}

		return otto.UndefinedValue()
	}
}*/

/*func privateDockerBuildCall(vm *notto.Notto, client *dockerclient.Client, options otto.Object) otto.Value {

}*/
