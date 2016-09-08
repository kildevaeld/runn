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
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all bundles",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		run, err := getRunn()
		if err != nil {
			printError(err)
		}

		bundles := run.List()
		if len(bundles) == 0 {
			fmt.Printf("No bundles")
			return
		}

		fmt.Printf("Bundles\n\n")
		/*ui.PaginatedList("", func(page int) []string {

			return nil
		})*/
		for _, bundle := range bundles {
			desc := bundle.Description
			if desc == "" {
				desc = "No description..."
			}
			fmt.Printf("%s\t\t%s\n", bundle.Name, bundle.Description)
		}

		fmt.Println("")
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.Aliases = []string{"ls"}
}
