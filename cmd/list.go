// Copyright © 2018 Ryan Parman <https://ryanparman.com>
// Copyright © 2018 Contributors <https://github.com/skyzyx/pstore/graphs/contributors>
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
	"os"
	"regexp"
	"sort"
	"strings"

	// "github.com/davecgh/go-spew/spew"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var parameters [][]string

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List values stored in Parameter Store.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		for _, page := range response {
			for _, param := range page.Parameters {
				parameters = append(parameters, []string{
					*param.Name,
					*param.Value,
				})
			}
		}

		if filter != "" {
			parameters = arrayFilter(parameters, func(v []string) bool {
				e := strings.Join(v, " ")
				return contains(e, filter)
			})
		} else if regex != "" {
			parameters = arrayFilter(parameters, func(v []string) bool {
				e := strings.Join([]string{v[0], v[1]}, " ")
				r, err := regexp.Compile(regex)

				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				return r.MatchString(e)
			})
		}

		if len(parameters) > 0 {
			// Sort alphabetically by key
			sort.Slice(parameters[:], func(i, j int) bool {
				return parameters[i][0] < parameters[j][0]
			})

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Key", "Value"})
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.AppendBulk(parameters)

			// Send output
			table.Render()
			fmt.Println("")
		}

		// Display result count
		results()
		fmt.Println("")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func results() {
	fmt.Printf("%s matched.\n", plural(len(parameters), "result", "results"))
}
