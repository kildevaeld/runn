package docker

import (
	"errors"
	"strings"

	dockerclient "github.com/fsouza/go-dockerclient"
	"github.com/kildevaeld/notto"
	"github.com/kildevaeld/notto/loop"
	"github.com/robertkrimen/otto"
)

func mustValue(value otto.Value, err error) otto.Value {
	if err != nil {
		panic(err)
	}
	return value
}

type docker_p_task struct {
	id       int64
	err      error
	result   interface{}
	callback otto.Value
	name     string
}

func (self *docker_p_task) GetID() int64   { return self.id }
func (self *docker_p_task) SetID(id int64) { self.id = id }

func (self *docker_p_task) Execute(vm *otto.Otto, l *loop.Loop) error {

	var arguments []interface{}

	if self.err != nil {
		e := vm.MakeCustomError("DockerError", self.err.Error())

		arguments = append(arguments, e)
	} else {
		arguments = append(arguments, otto.NullValue())
	}
	if self.result != nil {
		if v, e := vm.ToValue(self.result); e != nil {
			return e
		} else {
			arguments = append(arguments, v)
		}

	}

	if _, err := self.callback.Call(otto.NullValue(), arguments...); err != nil {
		return err
	}

	return nil

}

func (self *docker_p_task) Cancel() {

}

type docker_p struct {
	client *dockerclient.Client
	vm     *notto.Notto
}

func (self *docker_p) Build() {

}

func (self *docker_p) Create() {

}

func (self *docker_p) check_args(call otto.FunctionCall) error {
	if !call.Argument(0).IsObject() || call.Argument(0).Class() != "Object" {
		return errors.New("must " + call.Argument(0).Class())
	}
	if !call.Argument(1).IsFunction() {
		return errors.New("function")
	}
	return nil
}

func (self *docker_p) Start(call otto.FunctionCall) otto.Value {

	task := self.getTask("start", call.Argument(1))

	go func() {
		defer self.vm.Runloop().Ready(task)
		if task.err = self.check_args(call); task.err != nil {
			return
		}

		name := mustValue(call.Argument(0).Object().Get("name")).String()
		task.err = self.client.StartContainer(name, &dockerclient.HostConfig{})
	}()

	return otto.UndefinedValue()
}

func (self *docker_p) getTask(name string, call otto.Value) *docker_p_task {
	task := &docker_p_task{
		callback: call,
		name:     name,
	}

	self.vm.Runloop().Add(task)
	return task
}

func (self *docker_p) Stop(call otto.FunctionCall) otto.Value {
	task := self.getTask("stop", call.Argument(1))

	go func() {
		defer self.vm.Runloop().Ready(task)
		if task.err = self.check_args(call); task.err != nil {
			return
		}

		name := mustValue(call.Argument(0).Object().Get("name")).String()
		task.err = self.client.StopContainer(name, 10)

	}()

	return otto.UndefinedValue()
}

func (self *docker_p) HasImage(call otto.FunctionCall) otto.Value {
	task := self.getTask("has_image", call.Argument(1))

	go func() {
		defer self.vm.Runloop().Ready(task)
		if task.err = self.check_args(call); task.err != nil {
			return
		}

		name := mustValue(call.Argument(0).Object().Get("name")).String()
		images, err := self.client.ListImages(dockerclient.ListImagesOptions{})
		if err != nil {
			task.err = err
			return
		}

		for _, i := range images {
			if name == i.ID {
				task.result = true
				return
			}
			for _, t := range i.RepoTags {
				index := strings.Index(t, ":")

				if t[:index] == name {
					task.result = true
					return
				}
			}
		}

		task.result = false

	}()

	return otto.UndefinedValue()
}

func (self *docker_p) HasContainer() {

}

func (self *docker_p) RemoveContainer() {

}

func (self *docker_p) RemoveImage() {

}

func createDocker(vm *notto.Notto) (*docker_p, error) {
	c, e := dockerclient.NewClientFromEnv()
	if e != nil {
		return nil, e
	}
	return &docker_p{
		vm:     vm,
		client: c,
	}, nil
}
