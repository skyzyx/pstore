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
	"strconv"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

type Parameter struct {
	Profile string
	Name    string
	Type    string
	Value   string
}

const tmpl = `aws ssm put-parameter \
    --profile {{ .Profile }} \
    --name "{{ .Name }}" \
    --type {{ .Type }} \
    --value "{{ .Value }}" \
    --overwrite \
;
`

// cliCmd represents the cli command
var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Export the selected Parameter Store values as AWS CLI commands.",
	Run: func(cmd *cobra.Command, args []string) {
		if !quiet && !debug {
			fmt.Println("")
		}

		// Fetch the data from AWS
		cfg := GetConfig(cmd.Flag("profile").Value.String())
		svc := ssm.New(cfg)
		SendPathRequest(svc, args)

		for _, page := range response {
			for _, param := range page.Parameters {
				parameters = append(parameters, []string{
					*param.Name,
					*param.Value,
					string(param.Type),
					strconv.FormatInt(*param.Version, 10),
				})
			}
		}

		if filter != "" {
			parameters = ArrayFilter(parameters, func(v []string) bool {
				searchString := []string{v[0], v[1]}
				e := strings.Join(searchString, " ")

				if debug {
					colorize("Filterable string:")
					colorize(spew.Sdump(e))
					fmt.Println("")
				}

				return Contains(e, filter)
			})
		} else if regex != "" {
			parameters = ArrayFilter(parameters, func(v []string) bool {
				e := strings.Join([]string{v[0], v[1]}, " ")
				r, err := regexp.Compile(regex)

				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				if debug {
					colorize("Filterable string:")
					colorize(spew.Sdump(e, r))
					fmt.Println("")
				}

				return r.MatchString(e)
			})
		}

		if len(parameters) > 0 {
			// Sort alphabetically by key
			sort.Slice(parameters[:], func(i, j int) bool {
				return parameters[i][0] < parameters[j][0]
			})

			for _, entry := range parameters {
				construct := Parameter{
					Profile: profile,
					Name:    entry[0],
					Value:   entry[1],
					Type:    entry[2],
				}

				if debug {
					colorize("Parameter Object:")
					colorize(spew.Sdump(construct))
					fmt.Println("")
				}

				t := template.New("Parameter")

				t, err := t.Parse(tmpl)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				err = t.Execute(os.Stdout, construct)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				fmt.Println("")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(cliCmd)
}
