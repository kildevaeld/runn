package file

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kildevaeld/dict"
	"github.com/kildevaeld/runn/runnlib"
	"github.com/mitchellh/mapstructure"
)

type Config struct {
	Path string
}

type filestore struct {
	config Config
}

func (self *filestore) init() error {
	os.MkdirAll(filepath.Join(self.config.Path, "bundles"), 0766)
	return nil
}

func (self *filestore) Set(name string, r io.Reader, bundle runnlib.Bundle, length int64) error {
	path := filepath.Join(self.config.Path, name+".zip")
	file, e := os.Create(path)
	if e != nil {
		return e
	}
	defer file.Close()
	_, err := io.Copy(file, r)
	if err == nil {
		bundlePath := filepath.Join(self.config.Path, "bundles", name+".json")
		if b, e := json.MarshalIndent(bundle, "", "  "); e == nil {
			err = ioutil.WriteFile(bundlePath, b, 0666)
		} else {
			err = e
		}
	}

	return err

}

func (self *filestore) Get(name string) (io.Reader, error) {
	path := filepath.Join(self.config.Path, name+".zip")
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}

	return f, e

}

func (self *filestore) Remove(name string) error {
	path := filepath.Join(self.config.Path, name+".zip")
	bundlePath := filepath.Join(self.config.Path, "bundles", name+".json")
	os.RemoveAll(path)
	return os.RemoveAll(bundlePath)
}

func (self *filestore) List() []runnlib.Bundle {
	var out []runnlib.Bundle
	path := filepath.Join(self.config.Path, "bundles")
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("fs:list: %s\n", err)
		return out
	}
	//fmt.Printf("file store list")
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {

			continue
		}

		var b []byte
		if b, err = ioutil.ReadFile(filepath.Join(path, file.Name())); err != nil {
			fmt.Printf("fs:list: %s\n", err)
			return out
		}
		var bundle runnlib.Bundle
		if err = json.Unmarshal(b, &bundle); err != nil {
			fmt.Printf("fs:list: %s\n", err)
			return out
		}

		out = append(out, bundle)
	}

	return out
}

func init() {
	runnlib.AddStore("file", func(conf interface{}) (store runnlib.Store, err error) {
		var config Config
		switch t := conf.(type) {
		case Config:
			config = t
		case map[string]interface{}:
			err = mapstructure.Decode(t, &config)
		case dict.Map:
			err = mapstructure.Decode(t.ToMap(), &config)
		}
		if err == nil {
			fs := &filestore{config}
			err = fs.init()
			store = fs
		}

		return
	})
}
