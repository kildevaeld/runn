package cmd

import (
	"fmt"
	"os"

	"github.com/kildevaeld/runn"
	_ "github.com/kildevaeld/runn/store/file"
	_ "github.com/kildevaeld/runn/store/s3"
	"github.com/spf13/viper"
)

func printError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func getRunn() (*runn.Runn, error) {

	var config runn.Config

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("Config: %s", err)
	}

	//fmt.Printf("%#v", viper.AllKeys())

	return runn.New(config)
}
