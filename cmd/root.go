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

// The core functions which support this CLI tool.
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/caarlos0/spin"
	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var colorize func(format string, a ...interface{})
var debug bool
var filter string
var parameters [][]string
var profile string
var quiet bool
var regex string
var response []*ssm.GetParametersByPathOutput
var responseSingle *ssm.GetParameterOutput

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pstore",
	Short: "AWS Parameter Store Manager",
	Long: `Simplifies working with Parameter Store via the AWS CLI and Terraform.

Leverages the official AWS SDK under the hood, which means that all of the standard AWS CLI credential files and
environment variables can be used to configure this tool. The default behavior is to allow the AWS SDK to determine
the value from the environment.

Regular expression syntax can be found at https://github.com/google/re2/wiki/Syntax.`,
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
	colorize = color.New(color.FgRed).PrintfFunc()

	rootCmd.Version = "2.0.0"

	rootCmd.PersistentFlags().StringVar(&profile, "profile", "default",
		"(Optional) The AWS CLI Profile to use for the request.")

	rootCmd.PersistentFlags().StringVarP(&filter, "filter", "f", "",
		"(Optional) After the Parameter Store API call returns results, filter the names and values "+
			"by substring match.")

	rootCmd.PersistentFlags().StringVarP(&regex, "regex", "r", "",
		"(Optional) After the Parameter Store API call returns results, filter the names and values "+
			"by RE2 regular expression.")

	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false,
		"(Optional) Enable DEBUG logging.")

	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false,
		"(Optional) Do not display any messages during the fetching of data.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}

// Using the SDK’s default configuration, loading additional config
// and credentials values from the environment variables, shared
// credentials, and shared configuration files.
func GetConfig(profile string) aws.Config {
	if debug {
		colorize("Configuration profile:")
		colorize(spew.Sdump(profile))
		fmt.Println("")
	}

	cfg, err := external.LoadDefaultAWSConfig(external.WithSharedConfigProfile(profile))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return cfg
}

// Send the request for a PS Path to AWS, then stash the response into a global variable.
func SendPathRequest(svc *ssm.SSM, args []string) {
	path := "/"
	if len(args) > 0 {
		path = args[0]
	}

	if debug {
		colorize("Arguments:")
		colorize(spew.Sdump(args))
		fmt.Println("")
	}

	s := spin.New("Fetching %s ")
	if !quiet && !debug {
		s.Set(spin.Box2)
		s.Start()
		defer s.Stop()
	}

	input := &ssm.GetParametersByPathInput{
		MaxResults:     aws.Int64(10),
		Path:           &path,
		Recursive:      aws.Bool(true),
		WithDecryption: aws.Bool(true),
	}

	if debug {
		colorize("GetParametersByPathInput:")
		colorize(spew.Sdump(input))
		fmt.Println("")
	}

	// Create the request object without executing it
	request := svc.GetParametersByPathRequest(input)

	if debug {
		colorize("GetParametersByPathRequest:")
		colorize(spew.Sdump(request.Request.Config.Region))
		colorize(spew.Sdump(request.Request.Config.Credentials))
		fmt.Println("")
	}

	// Paginate the response
	p := request.Paginate()
	for p.Next() {
		if debug {
			colorize("Results Page:")
			colorize(spew.Sdump(p.CurrentPage()))
			fmt.Println("")
		}

		response = append(response, p.CurrentPage())
	}

	if err := p.Err(); err != nil {
		s.Stop()
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// Send the request for a single PS value to AWS, then stash the response into a global variable.
func SendSingleRequest(svc *ssm.SSM, args []string) {
	path := ""
	if len(args) > 0 {
		path = args[0]
	} else {
		fmt.Println("The path to the Parameter Store key is mandatory. If you don't know the path, try `pstore list` first.")
		os.Exit(1)
	}

	if debug {
		colorize("Arguments:")
		colorize(spew.Sdump(args))
		fmt.Println("")
	}

	s := spin.New("Fetching %s ")
	if !quiet && !debug {
		s.Set(spin.Box2)
		s.Start()
		defer s.Stop()
	}

	input := &ssm.GetParameterInput{
		Name:           &path,
		WithDecryption: aws.Bool(true),
	}

	if debug {
		colorize("GetParameterInput:")
		colorize(spew.Sdump(input))
		fmt.Println("")
	}

	// Create the request object without executing it
	request := svc.GetParameterRequest(input)

	if debug {
		colorize("GetParameterRequest:")
		colorize(spew.Sdump(request.Request.Config.Region))
		colorize(spew.Sdump(request.Request.Config.Credentials))
		fmt.Println("")
	}

	// Fetch the response
	var err error
	responseSingle, err = request.Send()

	if err != nil {
		s.Stop()
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if debug {
		colorize("Results:")
		colorize(spew.Sdump(responseSingle))
		fmt.Println("")
	}
}

// Pluralize a noun based on its count.
func Plural(count int, singular string, plural string) string {
	if count == 1 {
		return fmt.Sprintf("%d %s", count, singular)
	}

	return fmt.Sprintf("%d %s", count, plural)
}

// Case-insensitive strings.Contains().
func Contains(a, b string) bool {
	return strings.Contains(strings.ToUpper(a), strings.ToUpper(b))
}

// Returns a new slice containing all strings in the slice that satisfy
// the predicate f.
func ArrayFilter(vs [][]string, f func([]string) bool) [][]string {
	vsf := make([][]string, 0)

	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}

	return vsf
}
