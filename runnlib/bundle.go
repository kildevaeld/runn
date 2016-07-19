package runnlib

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kildevaeld/go-filecrypt"

	"gopkg.in/yaml.v2"
)

var NotExistsError = errors.New("NOENTRY")

type BundleCommand struct {
	Name    string
	Command CommandConfig
}

type Bundle struct {
	Name     string
	Commands []BundleCommand
	Context  []string
}

func GetBundleFromPath(dir string, v *Bundle) (err error) {
	p := filepath.Join(dir, "bundle.json")

	if _, err = os.Stat(p); err == nil {

		b, e := ioutil.ReadFile(p)
		if e != nil {
			return e
		}

		return json.Unmarshal(b, v)

	}

	p = filepath.Join(dir, "bundle.yaml")
	if _, err = os.Stat(p); err == nil {
		b, e := ioutil.ReadFile(p)
		if e != nil {
			return e
		}

		return yaml.Unmarshal(b, v)
	}

	return NotExistsError
}

func getContextFromDir(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var out []string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "bundle.") {
			continue
		}
		out = append(out, file.Name())
	}

	return out, nil

}

func getInterpreterFromExt(ext string) []string {
	switch ext {
	case ".sh":
		return []string{"sh", "-c"}
	case ".js":
		return []string{"node"}
	case ".py":
		return []string{"python"}
	default:
		return nil
	}
}

func PackageToDir(reader io.ReaderAt, size int64, target string, key_raw []byte) error {

	ar, err := zip.NewReader(reader, size)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(target, 0766); err != nil {
		return err
	}

	if st, err := os.Stat(target); err != nil {
		if err == os.ErrNotExist {
			if err = os.MkdirAll(target, 0666); err != nil {
				return err
			}
		}

	} else if !st.IsDir() {
		return errors.New("target is not a directory")
	}

	key := filecrypt.Key(key_raw)

	for _, file := range ar.File {
		path := filepath.Join(target, file.Name)
		f, err := os.Create(path)

		if err != nil {
			return fmt.Errorf("ERROR %v", err)
		}
		f.Chmod(0777)
		defer f.Close()

		r, e := file.Open()
		if e != nil {
			return e
		}
		defer r.Close()

		if e := filecrypt.Decrypt(f, r, &key); err != nil {
			return e
		}

	}

	return nil
}

func PackageFromDir(dir, name string, key_raw []byte) (*bytes.Buffer, error) {

	var bundle Bundle

	err := GetBundleFromPath(dir, &bundle)
	if err != nil && err != NotExistsError {
		return nil, err
	}

	if bundle.Name == "" {
		bundle.Name = name
	}

	w := bytes.NewBuffer(nil)

	ar := zip.NewWriter(w)

	if bundle.Context == nil {
		bundle.Context, err = getContextFromDir(dir)
		if err != nil {
			return nil, err
		}

		for _, b := range bundle.Context {
			base := b[0 : len(b)-len(filepath.Ext(b))]
			bundle.Commands = append(bundle.Commands, BundleCommand{
				Name: base,
				Command: CommandConfig{
					WorkDir:     fmt.Sprintf("{{.WorkDir}}"),
					Cmd:         fmt.Sprintf("{{.WorkDir}}/%s", b),
					Stdout:      "stdout",
					Stderr:      "stderr",
					Interpreter: getInterpreterFromExt(filepath.Ext(b)),
				},
			})
		}

	}

	key := filecrypt.Key(key_raw)

	for _, path := range bundle.Context {
		out, err := ar.Create(path)
		if err != nil {
			return nil, err
		}

		read, rerr := os.Open(filepath.Join(dir, path))
		if rerr != nil {
			return nil, rerr
		}

		if _, err := filecrypt.Encrypt(out, read, &key); err != nil {
			read.Close()
			return nil, err
		}

		/*if _, err := io.Copy(out, read); err != nil {
			read.Close()
			return nil, err
		}*/

		read.Close()
	}

	out, ferr := ar.Create("bundle.yaml")
	if ferr != nil {
		return nil, ferr
	}

	b, berr := yaml.Marshal(bundle)
	if berr != nil {
		return nil, berr
	}

	//out.Write(b)
	buf := bytes.NewBuffer(b)

	if _, err := filecrypt.Encrypt(out, buf, &key); err != nil {
		return nil, err
	}

	if err := ar.Close(); err != nil {
		return nil, err
	}

	return w, nil
}
