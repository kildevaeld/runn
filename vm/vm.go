package vm

import (
	"io"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/kildevaeld/notto"
	"github.com/kildevaeld/notto/modules/fetch"
	"github.com/kildevaeld/notto/modules/fs"
	"github.com/kildevaeld/notto/modules/global"
	"github.com/kildevaeld/notto/modules/process"
	"github.com/kildevaeld/notto/modules/promise"
	"github.com/kildevaeld/notto/modules/shell"
	"github.com/kildevaeld/notto/modules/ui"
	"github.com/kildevaeld/notto/modules/util"
	"github.com/kildevaeld/runn/vm/modules/archive"
	"github.com/kildevaeld/runn/vm/modules/docker"
	"github.com/kildevaeld/runn/vm/modules/s3"
)

func mustError(result error) func(err error) error {
	return func(err error) error {
		if err != nil {
			return multierror.Append(result, err)
		}
		return nil
	}
}

func NewVM(stdout, stderr io.Writer, workdir string, args []string, env map[string]string) (*notto.Notto, error) {

	vm := notto.New()

	a := vm.ProcessAttr()
	a.Argv = args
	a.Stderr = stderr
	a.Stdout = stdout
	a.Environ = notto.MapToEnviron(env)
	a.Cwd = workdir

	var result error
	result = mustError(result)(shell.Define(vm, false))
	result = mustError(result)(process.Define(vm))
	result = mustError(result)(util.Define(vm))
	result = mustError(result)(promise.Define(vm))
	result = mustError(result)(fetch.Define(vm))
	result = mustError(result)(global.Define(vm))
	result = mustError(result)(fs.Define(vm))
	result = mustError(result)(ui.Define(vm))

	result = mustError(result)(s3.Define(vm))
	result = mustError(result)(archive.Define(vm))
	result = mustError(result)(docker.Define(vm))

	if result != nil {
		return nil, result
	}

	return vm, nil
}
