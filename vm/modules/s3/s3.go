//go:generate go-bindata -pkg s3 -o s3_impl.go s3.js
package s3

import (
	"io"
	"os"
	"path/filepath"

	"github.com/kildevaeld/notto"
	"github.com/kildevaeld/notto/loop"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"github.com/robertkrimen/otto"
)

type s3Task struct {
	id           int64
	jsReq, jsRes *otto.Object
	cb           otto.Value
	err          error
	content      interface{}
	//status       int
	//statusText   string
	//headers      map[string][]string
	//body         []byte
}

func (t *s3Task) SetID(id int64) { t.id = id }
func (t *s3Task) GetID() int64   { return t.id }

func (t *s3Task) Execute(vm *otto.Otto, l *loop.Loop) error {
	var arguments []interface{}

	if t.err != nil {
		e, err := vm.Call(`new Error`, nil, t.err.Error())
		if err != nil {
			return err
		}

		arguments = append(arguments, e)
	} else {
		arguments = append(arguments, otto.NullValue())
	}
	if t.content != nil {
		if v, e := vm.ToValue(t.content); e != nil {
			return e
		} else {
			arguments = append(arguments, v)
		}

	}

	if _, err := t.cb.Call(otto.NullValue(), arguments...); err != nil {
		return err
	}

	return nil
}

func (t *s3Task) Cancel() {
}

func valToString(v otto.Value) string {
	if s, e := v.ToString(); e != nil {
		return ""
	} else {
		return s
	}
}

func mustValue(vm *notto.Notto) func(v otto.Value, err error) otto.Value {
	return func(v otto.Value, err error) otto.Value {
		if err != nil {
			vm.Throw("S3Error", err)
		}
		return v
	}
}

func mustStringValue(vm *notto.Notto) func(v string, err error) string {
	return func(v string, err error) string {
		if err != nil {
			vm.Throw("S3Error", err)
		}
		return v
	}
}

func Define(vm *notto.Notto) error {

	fn := func(call otto.FunctionCall) otto.Value {

		if len(call.ArgumentList) == 0 {
			vm.Throw("S3Error", "No options spcified")
		}

		if !call.Argument(0).IsObject() {
			vm.Throw("S3Error", "First argument must be an object")
		}

		o := call.Argument(0).Object()

		var (
			v    otto.Value
			e    error
			auth aws.Auth
		)

		if v, e = o.Get("access_key"); e != nil {
			vm.Throw("S3Error", "No access_key")
		}
		auth.AccessKey = valToString(v)
		if v, e = o.Get("secret_key"); e != nil {
			vm.Throw("S3Error", "no secret_key")
		}
		auth.SecretKey = valToString(v)

		bucket := mustValue(vm)(o.Get("bucket")).String()

		s := &s3_impl{s3.New(auth, aws.EUWest), bucket, vm}

		if v, e = vm.ToValue(s); e != nil {
			vm.Throw("S3Error", e)
		}
		return v
	}

	vm.Set("__private_s3", fn)

	vm.AddModule("s3", notto.CreateLoaderFromSource(string(MustAsset("s3.js")), ""))

	return nil
}

type s3_impl struct {
	client *s3.S3
	bucket string
	vm     *notto.Notto
}

func (self *s3_impl) Get(call otto.FunctionCall) otto.Value {

	path := mustStringValue(self.vm)(call.Argument(0).ToString())
	target := mustStringValue(self.vm)(call.Argument(1).ToString())
	t := &s3Task{
		cb: call.Argument(2),
	}

	self.vm.Runloop().Add(t)
	target = filepath.Join(self.vm.ProcessAttr().Cwd, target)
	go func() {
		defer self.vm.Runloop().Ready(t)

		var (
			err    error
			reader io.Reader
			file   *os.File
		)

		if reader, err = self.client.Bucket(self.bucket).GetReader(path); err != nil {
			t.err = err
			return
		}

		if file, err = os.Create(target); err != nil {
			t.err = err
			return
		}
		defer file.Close()

		if _, err = io.Copy(file, reader); err != nil {
			t.err = err
			return
		}

	}()

	return otto.UndefinedValue()
}

func (self *s3_impl) List(call otto.FunctionCall) otto.Value {

	prefix := mustStringValue(self.vm)(call.Argument(0).ToString())

	t := &s3Task{
		cb: call.Argument(1),
	}

	self.vm.Runloop().Add(t)

	go func() {
		defer self.vm.Runloop().Ready(t)

		var (
			err error
			res *s3.ListResp
		)

		if res, err = self.client.Bucket(self.bucket).List(prefix, "/", "", 1000); err != nil {
			t.err = err
			return
		}

		t.content = res.Contents

	}()

	return otto.UndefinedValue()
}
