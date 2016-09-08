package runn

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kildevaeld/runn/runnlib"
)

type Config struct {
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
	store runnlib.Store
	key   []byte
}

func (self *Runn) AddFromDir(path string) error {
	/*buf, err := runnlib.PackageFromDir(path, name, self.key)
	if err != nil {
		return err
	}*/
	buf, bundle, size, err := runnlib.ArchieveDir(path, self.key)
	if err != nil {
		return err
	}
	return self.store.Set(bundle.Name, buf, bundle, size)
}

func (self *Runn) Run(name, cmd string, config RunConfig) error {

	reader, err := self.store.Get(name)
	if err != nil {
		return fmt.Errorf("Store: %s", err)
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

	if err := runnlib.UnarchiveToDir(target, reader, 0, self.key); err != nil {
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

	store, serr := runnlib.GetStore(config.Store, config.StoreOptions)

	if serr != nil {
		return nil, serr
	}

	run := &Runn{store, []byte(config.Key)}

	return run, nil
}
