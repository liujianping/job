// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var envs *map[string]string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "job",
	Short: "Job, make your short-term command as a long-term job",
	Run: func(cmd *cobra.Command, args []string) {
		Main(cmd, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("config", "f", "", "job config file path")
	rootCmd.Flags().StringP("name", "N", "", "job name")
	// rootCmd.Flags().BoolP("cmd-http", "", false, "job command use inner http client, unsupport yet!!!")
	// rootCmd.Flags().BoolP("cmd-grpc", "", false, "job command use inner grpc client, unsupport yet!!!")
	// rootCmd.Flags().StringP("command", "C", "", "job command path name")

	envs = rootCmd.Flags().StringToStringP("cmd-env", "e", map[string]string{}, "job command enviromental variables")
	rootCmd.Flags().IntP("cmd-retry", "r", 0, "job command retry times when failed")
	rootCmd.Flags().DurationP("cmd-timeout", "t", 0, "job command timeout duration")

	rootCmd.Flags().IntP("concurrent", "c", 1, "job concurrent numbers ")
	rootCmd.Flags().IntP("repeat-times", "n", 1, "job repeat times, 0 means forever")
	rootCmd.Flags().DurationP("repeat-interval", "i", 0*time.Second, "job repeat interval duration")
	rootCmd.Flags().StringP("schedule", "s", "", "job schedule in crontab format")
	rootCmd.Flags().DurationP("timeout", "T", 0, "job timeout duration")
	rootCmd.Flags().BoolP("guarantee", "G", false, "job guarantee mode enable ?")
	rootCmd.Flags().BoolP("report", "R", false, "job reporter enable ?")
	viper.BindPFlags(rootCmd.Flags())
}
