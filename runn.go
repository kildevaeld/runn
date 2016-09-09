package runn

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/kildevaeld/runn/runnlib"
	"github.com/mitchellh/go-homedir"
)

type Config struct {
	Root         string `json:"config_path,omitempty"`
	Store        string
	StoreOptions interface{}
	Key          string
}

type RunConfig struct {
	Environ []string
	Args    []string
	Stdout  io.Writer
	Stderr  io.Writer
	Locals  map[string]interface{}
}

type Runn struct {
	store       runnlib.Store
	config      Config
	cachePath   string
	runtimePath string
	//key   []byte
}

func (self *Runn) AddFromDir(path string) error {
	/*buf, err := runnlib.PackageFromDir(path, name, self.key)
	if err != nil {
		return err
	}*/
	buf, bundle, size, err := runnlib.ArchieveDir(path, nil)
	if err != nil {
		return err
	}
	return self.store.Set(bundle.Name, buf, bundle, size)
}

func (self *Runn) ClearCache() error {
	os.RemoveAll(self.cachePath)
	return os.MkdirAll(self.cachePath, 0755)
}

func (self *Runn) Run(name, cmd string, config RunConfig) error {

	var reader *os.File
	var size int64
	var err error
	cache := filepath.Join(self.cachePath, name+".zip")
	if stat, err := os.Stat(cache); err == nil {
		mtime := stat.ModTime()
		diff := time.Now().Sub(mtime)
		size = stat.Size()
		if (diff / time.Minute) >= 5 {

		}
		reader, err = os.Open(cache)
	} else {
		tmp, err := self.store.Get(name)

		if err != nil {
			return fmt.Errorf("Store: %s", err)
		}
		if reader, err = os.Create(cache); err != nil {
			tmp.Close()
			return err
		}

		size, err = io.Copy(reader, tmp)
		tmp.Close()
		reader.Close()
		if err != nil {

			return err
		}

		reader, _ = os.Open(cache)

	}

	defer reader.Close()

	/*close := func() {}

	if closer, ok := reader.(io.ReadCloser); ok {
		close = func() { closer.Close() }
	}

	defer close()***/

	tmp := os.TempDir()
	target := filepath.Join(tmp, "runn", name)

	if err = os.MkdirAll(target, 0755); err != nil {
		return err
	}

	defer os.RemoveAll(target)

	if err := runnlib.UnarchiveToDir(target, reader, size, nil); err != nil {
		return err
	}

	bundle, berr := NewBundle(target)
	if berr != nil {
		return berr
	}

	return bundle.Run(cmd, config)

}

func (self *Runn) List() []runnlib.Bundle {
	return self.store.List()
}

func (self *Runn) Remove(name string) error {
	return self.store.Remove(name)
}

func New(config Config) (*Runn, error) {

	if config.Store == "" {
		return nil, errors.New("config needs a store")
	}

	var (
		store runnlib.Store
		err   error
	)
	if config.Root == "" {
		var home string
		if home, err = homedir.Dir(); err != nil {
			return nil, err
		}

		config.Root = filepath.Join(home, ".runn")
	}

	if _, err = os.Stat(""); err != nil {
		if err = os.MkdirAll(config.Root, 0755); err != nil {
			return nil, err
		}
	}

	cachePath := filepath.Join(config.Root, "cache")
	runtimePath := filepath.Join(config.Root, "runtime")

	os.MkdirAll(cachePath, 0755)
	os.MkdirAll(runtimePath, 0755)

	if store, err = runnlib.GetStore(config.Store, config.StoreOptions); err != nil {
		return nil, err
	}

	run := &Runn{store, config, cachePath, runtimePath}

	return run, nil
}
