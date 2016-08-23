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
	"os"
	"strings"

	"github.com/kildevaeld/runn"
	"github.com/spf13/cobra"
)

var envFlag []string

func mergeStrinSlices(slices ...[]string) []string {
	var out []string
	for _, slice := range slices {
		out = append(out, slice...)
	}
	return out
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long:  ``,
	//Args:    cli.RequiresMinArgs(1),
	Aliases: []string{"r"},
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		run, err := getRunn()
		if err != nil {
			printError(err)
		}

		if len(args) == 0 {
			printError(errors.New("no"))
		}

		split := strings.Split(args[0], ":")
		if len(split) == 1 {
			printError(errors.New("usage: runn run <bundle:command>"))
		}

		var a []string
		if len(args) > 1 {
			a = args[1:]
		}

		conf := runn.RunConfig{
			Environ: mergeStrinSlices(os.Environ(), envFlag),
			Args:    a,
		}

		if err = run.Run(split[0], split[1], conf); err != nil {
			printError(err)
		}

	},
}

func init() {
	RootCmd.AddCommand(runCmd)

	runCmd.Flags().StringSliceVarP(&envFlag, "env", "e", nil, "env")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
