//go:generate go-bindata -pkg archive -o archive_impl.go archive.js
package archive

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/kildevaeld/notto"
	"github.com/kildevaeld/notto/loop"
	"github.com/robertkrimen/otto"
)

type archivetask struct {
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

func (t *archivetask) SetID(id int64) { t.id = id }
func (t *archivetask) GetID() int64   { return t.id }

func (t *archivetask) Execute(vm *otto.Otto, l *loop.Loop) error {
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

func (t *archivetask) Cancel() {
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
			vm.Throw("ArchiveError", err)
		}
		return v
	}
}

func mustStringValue(vm *notto.Notto) func(v string, err error) string {
	return func(v string, err error) string {
		if err != nil {
			vm.Throw("ArchiveError", err)
		}
		return v
	}
}

func Define(vm *notto.Notto) error {

	loader := notto.CreateLoaderFromSource(string(MustAsset("archive.js")), "")

	vm.AddModule("archive", loader)
	vm.Set("__private_archive", func(call otto.FunctionCall) otto.Value {
		m := call.Argument(0).String()
		o := call.Argument(1).Object()

		if m != "pack" && m != "unpack" {
			vm.Throw("ArchiveError", "first parameter must be pack or unpack")
		}

		t := &archivetask{
			cb: call.Argument(2),
		}

		format := mustValue(vm)(o.Get("format")).String()
		source := mustValue(vm)(o.Get("source")).String()
		target := mustValue(vm)(o.Get("target")).String()

		vm.Runloop().Add(t)
		target = filepath.Join(vm.ProcessAttr().Cwd, target)
		source = filepath.Join(vm.ProcessAttr().Cwd, source)
		go func() {
			defer vm.Runloop().Ready(t)

			if m == "unpack" {
				t.err = unpack(format, source, target)
			} else {
				t.err = pack(format, source, target)
			}

		}()

		return otto.UndefinedValue()
	})

	return nil
}

func unpack(format, source, target string) error {

	var (
		file *os.File
		err  error
	)

	if file, err = os.Open(source); err != nil {
		return err
	}
	defer file.Close()
	if _, err = os.Stat(target); err != nil {
		if err = os.MkdirAll(target, 0766); err != nil {
			return err
		}
	}

	if format == "tar" {
		return unpack_tar(file, target, false)
	} else if format == "targz" {
		return unpack_tar(file, target, true)
	} else if format == "zip" {
		return unpack_zip(file, target)
	}

	return nil
}

func unpack_tar(source io.ReadCloser, target string, gz bool) error {
	var (
		err    error
		header *tar.Header
	)
	if gz {
		if source, err = gzip.NewReader(source); err != nil {
			return err
		}
		defer source.Close()
	}

	tarReader := tar.NewReader(source)

	for {
		if header, err = tarReader.Next(); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		filename := header.Name
		switch header.Typeflag {
		case tar.TypeDir:
			if err = os.MkdirAll(filepath.Join(target, filename), 0777); err != nil {
				return err
			}
		case tar.TypeReg:
			writer, err := os.Create(filepath.Join(target, filename))

			if err != nil {
				return err
			}

			io.Copy(writer, tarReader)
			writer.Close()
			if err = os.Chmod(filepath.Join(target, filename), os.FileMode(header.Mode)); err != nil {

				return err
			}
		default:
			//
		}

	}

	return nil
}

func unpack_zip(source *os.File, target string) error {
	size := int64(0)
	if state, err := source.Stat(); err != nil {
		return nil
	} else {
		size = state.Size()
	}

	reader, err := zip.NewReader(source, size)
	if err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return err
}

func pack(format, source, target string) error {
	return nil
}
