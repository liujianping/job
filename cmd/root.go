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
	"sort"
	"time"
	"context"

	"github.com/liujianping/job/config"
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
	Run: func(cmd *cobra.Command, args []string) {
		jds := []*config.JD{}
		if len(viper.GetString("config")) > 0 {
			cfs, err := config.ParseJDs(viper.GetString("config"))
			if err != nil {
				fmt.Println("job failed:", err)
				os.Exit(1)
			}
			jds = cfs
		}

		options := []config.Option{}
		options = append(options, config.Name(viper.GetString("name")))
		options = append(options, config.CommandName(viper.GetString("cmd-name")))
		for k, v := range viper.GetStringMapString("cmd-env") {
			options = append(options, config.CommandEnv(k, v))
		}
		options = append(options, config.CommandTimeout(viper.GetDuration("cmd-timeout")))
		options = append(options, config.CommandRetry(viper.GetInt("cmd-retry")))

		options = append(options, config.CommandArgs(args...))
		options = append(options, config.Crontab(viper.GetString("crontab")))
		options = append(options, config.RepeatTimes(viper.GetInt("repeat-times")))
		options = append(options, config.RepeatInterval(viper.GetDuration("repeat-interval")))
		options = append(options, config.Timeout(viper.GetDuration("timeout")))
		options = append(options, config.Concurrent(viper.GetInt("concurrent")))
		options = append(options, config.Guarantee(viper.GetBool("guarantee")))
		
		if len(jds) == 0 {
			jd := config.CommandJD()
			if viper.GetBool("cmd-http") {
				jd = config.HttpJD()
			}
			jds = append(jds, jd)
		}
		for _, jd := range jds {
			for _, opt := range options {
				opt(jd)
			}
		}
		sort.Sort(config.JDs(jds))

		var reporter *exec.Reporter
		mainOptions := []routine.Opt{routine.Interrupts(routine.DefaultCancelInterruptors...)}
		if viper.GetBool("report") {
			n := viper.GetInt("repeat-times") * viper.GetInt("concurrent")
			reporter =  exec.NewReporter(n)
			beforeExit := routine.ExecutorFunc(func(ctx context.Context) error {
				reporter.Finalize()
				return nil
			})
			mainOptions = append(mainOptions, routine.BeforeExit(beforeExit))
		}

		err := routine.Main(
			context.TODO(),
			routine.ExecutorFunc(func(ctx context.Context) error {
				if reporter != nil {
					routine.Go(ctx, reporter)
				}
				jobs := make(map[string]chan error, len(jds))
				var main chan error
				for _, jd := range jds {
					if len(jd.Order.Precondition) == 0 {
						job := exec.NewJob(jd, reporter)
						ch := routine.Go(ctx, job)
						jobs[job.String()] = ch
						if jd.Report {
							main = ch
						}
						if jd.Order.Wait {
							<-ch
						}
					}					
				}	
				for _, jd := range jds {
					for _, pre := range jd.Order.Precondition {
						if ch, ok := jobs[pre]; ok {
							<-ch
						}
						job := exec.NewJob(jd, reporter)
						ch := routine.Go(ctx, job)
						jobs[job.String()] = ch
						if jd.Report {
							main = ch
						}
						if jd.Order.Wait {
							<-ch
						}
					}				
				}			
				return <- main
			}),
			mainOptions...)
		if err != nil {
			fmt.Println("job failed:", err)
		}
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
	rootCmd.Flags().StringP("config", "f", "", "job config file path")
	rootCmd.Flags().StringP("name", "N", "", "job name")
	rootCmd.Flags().BoolP("cmd-http", "", false, "job command use http client")
	rootCmd.Flags().BoolP("cmd-grpc", "", false, "job command use grpc client")

	rootCmd.Flags().StringP("cmd-name", "C", "", "job command path name")
	rootCmd.Flags().StringArrayP("cmd-arg", "a", []string{""}, "job command argments")
	rootCmd.Flags().StringToStringP("cmd-env", "e", map[string]string{}, "job command enviromental variables")
	rootCmd.Flags().IntP("cmd-retry", "r", 0, "job command retry times when failed")
	rootCmd.Flags().DurationP("cmd-timeout", "t", 0, "job command timeout duration")

	rootCmd.Flags().IntP("concurrent", "c", 1, "job concurrent numbers ")
	rootCmd.Flags().IntP("repeat-times", "n", 1, "job repeat times, 0 means forever")
	rootCmd.Flags().DurationP("repeat-interval", "i", 0*time.Second, "job repeat interval duration")
	rootCmd.Flags().StringP("crontab", "", "", "job schedule plan in crontab format")
	rootCmd.Flags().DurationP("timeout", "T", 0, "job timeout duration")
	rootCmd.Flags().BoolP("guarantee", "G", false, "job guarantee mode enable ?")
	rootCmd.Flags().BoolP("report", "R", false, "job reporter enable ?")
	viper.BindPFlags(rootCmd.Flags())
}
