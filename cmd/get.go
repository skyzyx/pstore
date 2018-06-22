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

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/spf13/cobra"
)

// getCmd represents the cli command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets the value of a single Parameter Store key.",
	Run: func(cmd *cobra.Command, args []string) {
		if !quiet && !debug {
			fmt.Println("")
		}

		// Fetch the data from AWS
		cfg := GetConfig(cmd.Flag("profile").Value.String())
		svc := ssm.New(cfg)
		SendSingleRequest(svc, args)

		fmt.Println(*responseSingle.Parameter.Value)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
