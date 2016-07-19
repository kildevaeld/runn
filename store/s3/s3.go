package s3

import (
	"io"

	"github.com/kildevaeld/dict"
	"github.com/kildevaeld/runn/runnlib"
	"github.com/minio/minio-go"
	"github.com/mitchellh/mapstructure"
)

type Config struct {
	AccessKey    string
	AccessSecret string
	Bucket       string
}

type filestore struct {
	config Config
	client *minio.Client
}

func (self *filestore) init() error {
	client, err := minio.New("s3.amazonaws.com", self.config.AccessKey, self.config.AccessSecret, true)
	if err != nil {
		return err
	}

	self.client = client

	return nil
}

func (self *filestore) Set(name string, r io.Reader) error {
	_, err := self.client.PutObject(self.config.Bucket, name, r, "application/zip")
	return err

}

func (self *filestore) Get(name string) (io.Reader, int64, error) {

	r, e := self.client.GetObject(self.config.Bucket, name)
	if e != nil {
		return nil, 0, e
	}

	var stat minio.ObjectInfo
	if stat, e = r.Stat(); e != nil {
		return nil, 0, e
	}

	return r, stat.Size, e
}

func init() {
	runnlib.AddStore("s3", func(conf interface{}) (store runnlib.Store, err error) {
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
			fs := &filestore{config: config}
			err = fs.init()
			store = fs
		}

		return
	})
}
