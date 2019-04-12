package cmd

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/liujianping/job/config"
	"github.com/liujianping/job/exec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/x-mod/errors"
	"github.com/x-mod/routine"
)

func CheckErr(err error) {
	if err != nil {
		fmt.Println("job failed: ", errors.CauseFrom(err))
		os.Exit(errors.ValueFrom(err))
	}	
}

func Main(cmd *cobra.Command, args []string) {
	if len(viper.GetString("config")) == 0 && len(args) == 0 {
		CheckErr(errors.New("command required"))
	}
	jds := []*config.JD{}
	if len(viper.GetString("config")) > 0 {
		cfs, err := config.ParseJDs(viper.GetString("config"))
		if err != nil {
			fmt.Println("job failed:", err)
			os.Exit(1)
		}
		jds = cfs
	} else {
		jd := config.CommandJD()
		jd.Command.Name = args[0]
		jds = append(jds, jd)
	}

	options := []config.Option{}
	options = append(options, config.Name(viper.GetString("name")))
	for k, v := range *envs {
		options = append(options, config.CommandEnv(k, v))
	}
	options = append(options, config.CommandTimeout(viper.GetDuration("cmd-timeout")))
	options = append(options, config.CommandRetry(viper.GetInt("cmd-retry")))
	options = append(options, config.CommandArgs(args[1:]...))
	options = append(options, config.Crontab(viper.GetString("crontab")))
	options = append(options, config.RepeatTimes(viper.GetInt("repeat-times")))
	options = append(options, config.RepeatInterval(viper.GetDuration("repeat-interval")))
	options = append(options, config.Timeout(viper.GetDuration("timeout")))
	options = append(options, config.Concurrent(viper.GetInt("concurrent")))
	options = append(options, config.Guarantee(viper.GetBool("guarantee")))

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
		reporter = exec.NewReporter(n)
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
			return <-main
		}),
		mainOptions...)
	if err != nil {
		fmt.Println("job failed:", err)
	}
	os.Exit(errors.ValueFrom(err))
}
