package cmd

import (
	"context"
	"fmt"
	"os"
	"sort"

	"gopkg.in/yaml.v2"

	"github.com/liujianping/job/config"
	"github.com/liujianping/job/exec"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/x-mod/errors"
	"github.com/x-mod/routine"
)

func CheckErr(err error) {
	if err != nil {
		fmt.Println("job failed: ", errors.CauseFrom(err))
		fmt.Println()
	}
}

func Main(cmd *cobra.Command, args []string) {
	if len(viper.GetString("config")) == 0 && len(args) == 0 {
		cmd.Usage()
		os.Exit(0)
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
		jd.Command.Shell.Name = args[0]
		jds = append(jds, jd)
	}

	options := []config.Option{}
	options = append(options, config.Name(viper.GetString("name")))
	for k, v := range *metadata {
		options = append(options, config.Metadata(k, v))
	}
	for k, v := range *envs {
		options = append(options, config.CommandEnv(k, v))
	}
	options = append(options, config.CommandTimeout(viper.GetDuration("cmd-timeout")))
	options = append(options, config.CommandRetry(viper.GetInt("cmd-retry")))
	if len(args) > 1 {
		options = append(options, config.CommandArgs(args[1:]...))
	}
	options = append(options, config.CommandStdoutDiscard(viper.GetBool("cmd-stdout-discard")))

	options = append(options, config.Guarantee(viper.GetBool("guarantee")))
	options = append(options, config.Crontab(viper.GetString("crontab")))
	options = append(options, config.RepeatTimes(viper.GetInt("repeat-times")))
	options = append(options, config.RepeatInterval(viper.GetDuration("repeat-interval")))
	options = append(options, config.Timeout(viper.GetDuration("timeout")))
	options = append(options, config.Concurrent(viper.GetInt("concurrent")))

	for _, jd := range jds {
		for _, opt := range options {
			opt(jd)
		}
	}
	sort.Sort(config.JDs(jds))

	//output
	if viper.GetBool("output") {
		for i, jd := range jds {
			bt, err := yaml.Marshal(map[string]*config.JD{
				"Job": jd,
			})
			if err != nil {
				fmt.Println("job failed:", err)
				os.Exit(errors.ValueFrom(err))
			}
			if i > 0 {
				fmt.Println("---")
			}
			fmt.Print(string(bt))
		}
		os.Exit(0)
	}

	var reporter *exec.Reporter
	mainOptions := []routine.Opt{routine.Interrupts(routine.DefaultCancelInterruptors...)}
	if viper.GetBool("report") {
		n := viper.GetInt("repeat-times") * viper.GetInt("concurrent")
		reporter = exec.NewReporter(n)
		beforeExit := routine.ExecutorFunc(func(ctx context.Context) error {
			reporter.Stop()
			reporter.Finalize()
			return nil
		})
		mainOptions = append(mainOptions, routine.BeforeExit(beforeExit))
	}

	ctx := context.TODO()
	if viper.GetBool("verbose") {
		log := logrus.New()
		log.SetLevel(logrus.TraceLevel)
		ctx = routine.WithLogger(ctx, log)
	}

	err := routine.Main(
		ctx,
		routine.ExecutorFunc(func(ctx context.Context) error {
			if reporter != nil {
				routine.Go(ctx, reporter)
			}
			jobs := make(map[string]chan error, len(jds))
			var main chan error
			var tail chan error
			for _, jd := range jds {
				if len(jd.Order.Precondition) == 0 {
					job := exec.NewJob(jd, reporter)
					ch := routine.Go(ctx, job)
					jobs[job.String()] = ch
					if jd.Report {
						main = ch
					}
					tail = ch
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
					tail = ch
					if jd.Order.Wait {
						<-ch
					}
				}
			}
			if main == nil {
				main = tail
			}
			return <-main
		}),
		mainOptions...)
	if err != nil {
		fmt.Println("job failed:", err)
	}
	os.Exit(errors.ValueFrom(err))
}
