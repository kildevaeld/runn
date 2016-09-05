package vm

import (
	"io"

	"github.com/kildevaeld/notto"
	"github.com/kildevaeld/notto/modules/fetch"
	"github.com/kildevaeld/notto/modules/global"
	"github.com/kildevaeld/notto/modules/process"
	"github.com/kildevaeld/notto/modules/promise"
	"github.com/kildevaeld/notto/modules/shell"
	"github.com/kildevaeld/notto/modules/util"
	"github.com/kildevaeld/runn/vm/modules/s3"
)

func NewVM(stdout, stderr io.Writer, workdir string, args []string, env map[string]string) (*notto.Notto, error) {

	vm := notto.New()

	a := vm.ProcessAttr()
	a.Argv = args
	a.Stderr = stderr
	a.Stdout = stdout
	a.Environ = notto.MapToEnviron(env)
	a.Cwd = workdir

	shell.Define(vm, false)
	process.Define(vm)
	util.Define(vm)
	promise.Define(vm)
	fetch.Define(vm)
	global.Define(vm)

	// Custom modules
	if e := s3.Define(vm); e != nil {
		panic(e)
	}

	/*ob, err := vm.Object("({})")
	if err != nil {
		return nil, err
	}

	getStringList := func(call otto.FunctionCall) []string {
		var out []string
		for _, a := range call.ArgumentList {
			if a.IsString() {
				if s, e := a.ToString(); e == nil {
					out = append(out, s)
				}
			} else if a.IsNumber() {
				if s, e := a.ToInteger(); e == nil {
					out = append(out, fmt.Sprintf("%d", s))
				}
			} else {
				continue
			}

		}
		return out
	}

	ob.Set("log", func(call otto.FunctionCall) otto.Value {
		str := getStringList(call)
		stdout.Write([]byte(strings.Join(str, " ")))
		return otto.UndefinedValue()
	})

	ob.Set("error", func(call otto.FunctionCall) otto.Value {
		str := getStringList(call)
		stderr.Write([]byte(strings.Join(str, " ")))
		return otto.UndefinedValue()
	})

	vm.Set("console", ob)*/

	return vm, nil
}
