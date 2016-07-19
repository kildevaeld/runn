package file

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

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
	os.MkdirAll(self.config.Path, 0766)
	return nil
}

func (self *filestore) Set(name string, r io.Reader) error {
	path := filepath.Join(self.config.Path, name+".zip")
	b, e := ioutil.ReadAll(r)
	if e != nil {
		return e
	}

	return ioutil.WriteFile(path, b, 0766)

}

func (self *filestore) Get(name string) (io.ReaderAt, int64, error) {
	path := filepath.Join(self.config.Path, name+".zip")
	f, e := os.Open(path)
	if e != nil {
		return nil, 0, e
	}
	var stat os.FileInfo
	if stat, e = f.Stat(); e != nil {
		return nil, 0, e
	}

	return f, stat.Size(), e

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
