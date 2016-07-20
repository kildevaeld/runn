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

type Runn struct {
	store runnlib.Store
	key   []byte
}

func (self *Runn) AddFromDir(name, path string) error {
	/*buf, err := runnlib.PackageFromDir(path, name, self.key)
	if err != nil {
		return err
	}*/
	buf, bundle, size, err := runnlib.ArchieveDir(path, name, self.key)
	if err != nil {
		return err
	}
	return self.store.Set(name, buf, bundle, size)
}

func (self *Runn) Run(name, cmd string) error {

	reader, err := self.store.Get(name)
	if err != nil {
		return fmt.Errorf("Store: %s", err)
	}

	close := func() {}

	if closer, ok := reader.(io.ReadCloser); ok {
		close = func() { closer.Close() }
	}

	defer close()

	tmp := os.TempDir()
	target := filepath.Join(tmp, "runn", name)

	if err = os.MkdirAll(target, 0766); err != nil {
		return err
	}

	defer os.RemoveAll(target)

	/*if err := runnlib.PackageToDir(reader, size, target, self.key); err != nil {
		return fmt.Errorf("PackageToDir: %s", err)
	}*/

	if err := runnlib.UnarchiveToDir(target, reader, 0, self.key); err != nil {
		return err
	}

	bundle, berr := NewBundle(target)
	if berr != nil {
		return berr
	}

	return bundle.Run(cmd)

}

func (self *Runn) List() []runnlib.Bundle {
	return self.store.List()
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
