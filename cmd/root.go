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
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/caarlos0/spin"
	"github.com/spf13/cobra"
)

var awsKey string
var filter string
var profile string
var regex string
var region string
var response []*ssm.GetParametersByPathOutput
var secretKey string
var token string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pstore",
	Short: "AWS Parameter Store Manager",
	Long: `Simplifies working with Parameter Store via the AWS CLI and Terraform.

Leverages the official AWS SDK under the hood, which means that all of the standard AWS CLI credential files and
environment variables can be used to configure this tool. The default behavior is to allow the AWS SDK to determine
the value from the environment.

Regular expression syntax can be found at https://github.com/google/re2/wiki/Syntax.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		// Don't execute for help
		if cmd.Use != "help [command]" {
			fmt.Println("")

			// Get AWS config for selected profile
			cfg := getConfig(cmd.Flag("profile").Value.String())
			svc := ssm.New(cfg)

			sendRequest(svc, args)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Initialize the shared flags that are valid across all commands.
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Version = "2.0.0"

	rootCmd.PersistentFlags().StringVar(&profile, "profile", "default",
		"(Optional) The AWS CLI Profile to use for the request.")

	rootCmd.PersistentFlags().StringVarP(&filter, "filter", "f", "",
		"(Optional) After the Parameter Store API call returns results, filter the names and values " +
		"by substring match.")

	rootCmd.PersistentFlags().StringVarP(&regex, "regex", "r", "",
		"(Optional) After the Parameter Store API call returns results, filter the names and values " +
		"by RE2 regular expression.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}

// Using the SDK's default configuration, loading additional config
// and credentials values from the environment variables, shared
// credentials, and shared configuration files
func getConfig(profile string) aws.Config {
	cfg, err := external.LoadDefaultAWSConfig(external.WithSharedConfigProfile(profile))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return cfg
}

// Send the request to AWS, then stash the response into a variable.
func sendRequest(svc *ssm.SSM, args []string) {
	path := "/"
	if len(args) > 0 {
		path = args[0]
	}

	s := spin.New("Fetching %s ")
	s.Set(spin.Box2)
	s.Start()
	defer s.Stop()

	// Create the request object without executing it
	request := svc.GetParametersByPathRequest(&ssm.GetParametersByPathInput{
		MaxResults:     aws.Int64(10),
		Path:           &path,
		Recursive:      aws.Bool(true),
		WithDecryption: aws.Bool(true),
	})

	// Paginate the response
	p := request.Paginate()
	for p.Next() {
		response = append(response, p.CurrentPage())
	}

	if err := p.Err(); err != nil {
		s.Stop()
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// Pluralize a noun based on its count
func plural(count int, singular string, plural string) string {
	if count == 1 {
		return fmt.Sprintf("%d %s", count, singular)
	}

	return fmt.Sprintf("%d %s", count, plural)
}

// Case-insensitive strings.Contains().
func contains(a, b string) bool {
	return strings.Contains(strings.ToUpper(a), strings.ToUpper(b))
}

// Returns a new slice containing all strings in the slice that satisfy
// the predicate `f`.
func arrayFilter(vs [][]string, f func([]string) bool) [][]string {
	vsf := make([][]string, 0)

	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}

	return vsf
}
