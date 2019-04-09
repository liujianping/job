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
	"strings"
	"context"
	"fmt"
	"os"
	"time"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/liujianping/job/exec"
	"github.com/x-mod/routine"
	"github.com/x-mod/errors"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "job",
	Short: "Job, make your short-term command as a long-term job",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		options := []exec.Option{}
		options = append(options, exec.JobName(viper.GetString("name")))
		options = append(options, exec.JobCommand(viper.GetString("cmd")))
		for k, v := range viper.GetStringMapString("env") {
			options = append(options, exec.JobEnv(k, v))
		}
		arguments := []string{}
		arguments = append(arguments, args...)	
		arguments = append(arguments, strings.Split(strings.TrimRight(strings.TrimLeft(viper.GetString("arg"), "["), "]"), ",")...)	
		options = append(options, exec.JobArgs(arguments...))
		options = append(options, exec.JobPayload([]byte(viper.GetString("payload"))))

		options = append(options, exec.JobCrontab(viper.GetString("crontab")))
		options = append(options, exec.JobRetryTimes(viper.GetInt32("retry-times")))
		options = append(options, exec.JobRepeatTimes(viper.GetInt32("repeat-times")))
		options = append(options, exec.JobRepeatInterval(viper.GetDuration("repeat-interval")))
		options = append(options, exec.CmdTimeout(viper.GetDuration("cmd-timeout")))
		options = append(options, exec.JobTimeout(viper.GetDuration("job-timeout")))
		// options = append(options, exec.JobTimeout(viper.GetDuration("deadline")))
		options = append(options, exec.PipeLineMode(viper.GetBool("pipeline-mode")))
		options = append(options, exec.Guarantee(viper.GetBool("guarantee")))
		options = append(options, exec.Concurrent(viper.GetInt("concurrent")))
		mainOptions := []routine.Opt{routine.Interrupts(routine.DefaultCancelInterruptors...)}
		if viper.GetBool("report") {
			results := make(chan *routine.Result, 10000)
			options = append(options, exec.Report(results))
			n := viper.GetInt("repeat-times") * viper.GetInt("concurrent")
			mainOptions = append(mainOptions, routine.BeforeExit(exec.NewReporter(n, results)))
		}
					
		err := routine.Main(
			context.TODO(), 
			exec.NewJob(options...),
			mainOptions...)
		os.Exit(errors.ValueFrom(err))
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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.job.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringP("name", "", "", "job name")
	rootCmd.Flags().StringP("cmd", "C", "", "job command path")
	rootCmd.Flags().StringArrayP("arg", "a", []string{""}, "job command argments")
	rootCmd.Flags().StringToStringP("env", "e", map[string]string{}, "job command enviromental variables")
	rootCmd.Flags().BytesBase64P("payload", "p", []byte{}, "job command custom payload")

	rootCmd.Flags().IntP("concurrent", "c", 1, "job command concurrent numbers ")
	rootCmd.Flags().IntP("retry-times", "r", 0, "job command retry times when failed")
	rootCmd.Flags().IntP("repeat-times", "n", 1, "job command repeat times, 0 means forever")
	rootCmd.Flags().DurationP("repeat-interval", "i", 0*time.Second, "job command repeat interval duration")
	rootCmd.Flags().StringP("crontab", "", "", "job command schedule plan in crontab format")
	rootCmd.Flags().DurationP("cmd-timeout", "t", 0, "timeout duration for command")
	rootCmd.Flags().DurationP("job-timeout", "T", 0, "timeout duration for the job")

	rootCmd.Flags().BoolP("pipline", "P", false, "job executing in pipeline mode")
	rootCmd.Flags().BoolP("guarantee", "G", false, "job executing in guarantee mode")
	rootCmd.Flags().BoolP("report", "R", false, "job generate report")
	viper.BindPFlags(rootCmd.Flags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".job" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".job")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
