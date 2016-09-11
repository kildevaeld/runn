package docker

import (
	"testing"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/kildevaeld/notto"
	"github.com/kildevaeld/notto/modules"
)

func mustError(result error) func(err error) error {
	return func(err error) error {
		if err != nil {
			return multierror.Append(result, err)
		}
		return nil
	}
}

func TBuilder(t *testing.T) {

	vm := notto.New()

	/*var result error
	result = mustError(result)(shell.Define(vm, false))
	result = mustError(result)(process.Define(vm))
	result = mustError(result)(util.Define(vm))
	result = mustError(result)(promise.Define(vm))
	result = mustError(result)(fetch.Define(vm))
	result = mustError(result)(global.Define(vm))
	result = mustError(result)(fs.Define(vm))
	result = mustError(result)(ui.Define(vm))
	result = mustError(result)(fsm.Define(vm))
	result = mustError(result)(s3.Define(vm))
	result = mustError(result)(archive.Define(vm))
	result = mustError(result)(Define(vm))*/
	result := modules.Define(vm)
	result = Define(vm)
	if result != nil {
		t.Fatal(result)
	}

	var s = `
        var docker = require('docker.builder');
        var config = require('./test');
        docker.createBuilder(config,'production')
        .then(function (builder) {

            builder.on('notification', function (e, m) {
                console.log(e,m.name)
            });

            return builder.start(true);
            /*builder.build()
            .then(function () {
                return builder.start();
            }).catch(console.log)*/
        })
        
    `

	_, e := vm.RunScript(s, "")
	if e != nil {
		t.Fatal(e)
	}
}

func TestDocker(t *testing.T) {

	vm := notto.New()

	result := modules.Define(vm)
	result = Define2(vm)
	if result != nil {
		t.Fatal(result)
	}

	s := `
		var docker = require('docker').create();

		docker.Start({name: "nginx"}, function (err) {
			console.log('started: ', err)
			docker.HasImage({name: 'sha256:1caaa506d2e83f60781476e5732a7b70e0443ab848cc484085fb2984b40fa805', timeout:1}, function (e, ok) {
				console.log(e, ok)
			});
		});
	
	`

	if _, e := vm.RunScript(s, ""); e != nil {
		t.Fatal(e)
	}

}
