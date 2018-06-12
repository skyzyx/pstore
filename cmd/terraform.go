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

	"github.com/spf13/cobra"
)

// terraformCmd represents the terraform command
var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Export the selected Parameter Store values as Terraform `*.tf` files.",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("terraform called")
	},
}

func init() {
	// rootCmd.AddCommand(terraformCmd)
}
