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
	"bytes"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/kildevaeld/runn/runnlib"
	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen <path>",
	Short: "Generate bundle from directory",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("gen called")

		if len(args) == 0 {
			printError(errors.New("no args"))
		}

		path, err := filepath.Abs(args[0])
		if err != nil {
			printError(err)
		}
		var buf io.Reader
		/*if buf, err = runnlib.PackageFromDir(path, "", []byte("tmp")); err != nil {
			printError(err)
		}*/
		var size int64
		if buf, _, size, err = runnlib.ArchieveDir(path, []byte("tmp")); err != nil {
			printError(err)
		}

		reader := buf.(*bytes.Buffer)

		if err := runnlib.UnarchiveToDir(path, bytes.NewReader(reader.Bytes()), size, []byte("tmp")); err != nil {
			printError(err)
		}
		//err = runnlib.PackageToDir(bytes.NewReader(reader.Bytes()), int64(buf.Len()), path, []byte("tmp"))
		//printError(err)
	},
}

func init() {
	RootCmd.AddCommand(genCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
