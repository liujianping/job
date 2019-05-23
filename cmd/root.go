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
	"fmt"
	"os"
	"time"

	"github.com/liujianping/job/build"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var envs *map[string]string
var metadata *map[string]string

var rootCmd *cobra.Command

//RootCmd new root cmd
func RootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "job [flags] [command args ...]",
		Short: "Job, make your short-term command as a long-term job",
		Example: `
	(simple)      $: job echo hello	
	(schedule)    $: job -s "* * * * *" -- echo hello
	(retry)       $: job -r 3 -- echox hello
	(repeat)      $: job -n 10 -i 100ms -- echo hello
	(concurrent)  $: job -c 10 -n 10 -- echo hello
	(report)      $: job -R -- echo hello
	(timeout cmd) $: job -t 500ms -R -- sleep 1
	(timeout job) $: job -n 10 -i 500ms -T 3s -R -- echo hello
	(job output)  $: job -n 10 -i 500ms -T 3s -o -- echo hello
	(job config)  $: job -f /path/to/job.yaml`,
		Run: func(cmd *cobra.Command, args []string) {
			if viper.GetBool("version") {
				fmt.Println(build.String())
				os.Exit(0)
			}

			if len(viper.GetString("config")) == 0 && len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
			exitForErr(Main(cmd, args))
		},
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd = RootCmd()
	rootCmd.Flags().StringP("config", "f", "", "job config file path")
	rootCmd.Flags().StringP("name", "N", "", "job name definition")
	metadata = rootCmd.Flags().StringToStringP("metadata", "M", map[string]string{}, "job metadata definition")
	envs = rootCmd.Flags().StringToStringP("cmd-env", "e", map[string]string{}, "job command environmental variables")
	rootCmd.Flags().IntP("cmd-retry", "r", 0, "job command retry times when failed")
	rootCmd.Flags().DurationP("cmd-timeout", "t", 0, "job command timeout duration")
	rootCmd.Flags().BoolP("cmd-stdout-discard", "d", false, "job command stdout discard ?")

	rootCmd.Flags().IntP("concurrent", "c", 0, "job concurrent numbers ")
	rootCmd.Flags().IntP("repeat-times", "n", 1, "job repeat times, 0 means forever")
	rootCmd.Flags().DurationP("repeat-interval", "i", 0*time.Second, "job repeat interval duration")
	rootCmd.Flags().StringP("schedule", "s", "", "job schedule in crontab format")
	rootCmd.Flags().DurationP("timeout", "T", 0, "job timeout duration")
	rootCmd.Flags().BoolP("guarantee", "G", false, "job guarantee mode enable ?")
	rootCmd.Flags().BoolP("report", "R", false, "job report enable ?")
	rootCmd.Flags().StringP("report-push-gateway", "P", "", "job report to prometheus push gateway address")
	rootCmd.Flags().DurationP("report-push-interval", "I", 0*time.Second, "job report to prometheus push gateway interval")
	rootCmd.Flags().BoolP("output", "o", false, "job yaml config output enable ?")
	// rootCmd.Flags().StringP("output-command-format", "F", "shell", "job yaml config output command format ?")
	rootCmd.Flags().BoolP("verbose", "V", false, "job verbose log enable ?")
	rootCmd.Flags().BoolP("version", "v", false, "job version")

	// TODO support Distributed-Job
	// rootCmd.Flags().StringP("host", "H", "", "dispatch JOB to the Host")
	viper.BindPFlags(rootCmd.Flags())
	rootCmd.HelpFunc()
}
