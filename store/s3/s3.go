package s3

import (
	"encoding/json"
	"io"
	"log"

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

type s3filestore struct {
	config Config
	bucket *s3.Bucket
}

func (self *s3filestore) init() error {

	auth := aws.Auth{
		AccessKey: self.config.AccessKey,
		SecretKey: self.config.AccessSecret,
	}

	client := s3.New(auth, aws.EUWest)

	self.bucket = client.Bucket(self.config.Bucket)

	return nil
}

func (self *s3filestore) Set(name string, r io.Reader, bundle runnlib.Bundle, length int64) error {
	//_, err := self.client.PutObject(self.config.Bucket, name, r, "application/zip")
	//return err
	err := self.bucket.PutReader(name, r, length, "application/zip", s3.Private)

	if err == nil {
		var b []byte
		if b, err = json.Marshal(bundle); err != nil {
			return err
		}
		return self.bucket.Put("bundles/"+name+".json", b, "application/json", s3.Private)
	}

	return err
}

func (self *s3filestore) Get(name string) (io.ReadCloser, error) {

	r, e := self.bucket.GetReader(name)

	//r, e := self.client.GetObject(self.config.Bucket, name)
	if e != nil {
		return nil, e
	}

	return r, e
}

func (self *s3filestore) List() []runnlib.Bundle {
	var out []runnlib.Bundle
	r, e := self.bucket.List("bundles/", "/", "", 1000)
	if e != nil {
		return out
	}

	if r.Contents != nil {
		for _, k := range r.Contents {
			//key := strings.TrimPrefix(k.Key, "bundles/")
			b, e := self.bucket.Get(k.Key)
			if e != nil {
				log.Printf("s3: list: %s\n", e)
				continue
			}

			var bundle runnlib.Bundle
			if err := json.Unmarshal(b, &bundle); err != nil {
				continue
			}
			out = append(out, bundle)
		}
	}

	return out
}

func (self *s3filestore) Remove(name string) error {
	e := self.bucket.Del(name)
	if e != nil {
		return e
	}
	return self.bucket.Del("bundles/" + name + ".json")
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
			fs := &s3filestore{config: config}
			err = fs.init()
			store = fs
		}

		return
	})
}
