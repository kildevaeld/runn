//go:generate go-bindata -pkg docker -o docker_impl.go docker.js
package docker

import "github.com/kildevaeld/notto"

func Define(vm *notto.Notto) error {

	//err := vm.Set("__private_docker", privateDockerCall(vm))

	vm.AddModule("docker", notto.CreateLoaderFromSource(string(MustAsset("docker.js")), ""))

	return nil
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
