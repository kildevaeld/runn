// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"
	"os"
	"text/template"

	"github.com/kildevaeld/runn/runnlib"
	"github.com/spf13/cobra"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:   "describe <bundle>",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			printError(errors.New("usage: runn describe <bundle>"))
		}

		run, err := getRunn()
		if err != nil {
			printError(err)
		}

		var bundle runnlib.Bundle
		found := false
		for _, b := range run.List() {
			if b.Name == args[0] {
				bundle = b
				found = true
				break
			}
		}

		if !found {
			fmt.Printf("Bundle with name '%s' not found", args[0])
		}

		templ, e := template.New("description").Parse(string(MustAsset("templates/bundle.tmpl")))
		if e != nil {
			printError(e)
		}
		templ.Execute(os.Stdout, bundle)
	},
}

func init() {
	RootCmd.AddCommand(describeCmd)
	describeCmd.Aliases = []string{"desc"}

}
