package s3

import (
	"io"

	"github.com/kildevaeld/dict"
	"github.com/kildevaeld/runn/runnlib"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"github.com/mitchellh/mapstructure"
)

type Config struct {
	AccessKey    string
	AccessSecret string
	Bucket       string
}

type filestore struct {
	config Config
	bucket *s3.Bucket
}

func (self *filestore) init() error {

	auth := aws.Auth{
		AccessKey: self.config.AccessKey,
		SecretKey: self.config.AccessSecret,
	}

	client := s3.New(auth, aws.EUWest)

	self.bucket = client.Bucket(self.config.Bucket)

	return nil
}

func (self *filestore) Set(name string, r io.Reader, bundle runnlib.Bundle, length int64) error {
	//_, err := self.client.PutObject(self.config.Bucket, name, r, "application/zip")
	//return err
	return self.bucket.PutReader(name, r, length, "application/zip", s3.Private)
}

func (self *filestore) Get(name string) (io.Reader, error) {

	r, e := self.bucket.GetReader(name)

	//r, e := self.client.GetObject(self.config.Bucket, name)
	if e != nil {
		return nil, e
	}

	return r, e
}

func (self *filestore) List() []runnlib.Bundle {

	return nil
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
