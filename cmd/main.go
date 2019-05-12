package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/liujianping/job/config"
	"github.com/liujianping/job/exec"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/x-mod/errors"
	"github.com/x-mod/httpclient"
	"github.com/x-mod/routine"
	"gopkg.in/yaml.v2"
)

var (
	needReport      = false
	httpConnections = 0
)

func optionJDs(jds []*config.JD, options []config.Option, report bool) {
	for _, jd := range jds {
		for _, opt := range options {
			opt(jd)
			if report {
				jd.Report = true
			}
			if jd.Report {
				needReport = true
			}
			if jd.Command.HTTP != nil {
				httpConnections = httpConnections + jd.Concurrent
			}
		}
	}
}

func output(jds []*config.JD) {
	for i, jd := range jds {
		bt, err := yaml.Marshal(map[string]*config.JD{
			"Job": jd,
		})
		exitForErr(err)
		if i > 0 {
			fmt.Println("---")
		}
		fmt.Print(string(bt))
	}
}

func exitForErr(err error) {
	if err != nil {
		fmt.Println("job failed:", err)
		os.Exit(errors.ValueFrom(err))
	}
}

func withVerbose(ctx context.Context) context.Context {
	if viper.GetBool("verbose") {
		log := logrus.New()
		log.SetLevel(logrus.TraceLevel)
		return routine.WithLogger(ctx, log)
	}
	return ctx
}

//Main func
func Main(cmd *cobra.Command, args []string) {
	jds := []*config.JD{}
	if len(viper.GetString("config")) > 0 {
		cfs, err := config.ParseJDs(viper.GetString("config"))
		exitForErr(err)
		jds = append(jds, cfs...)
	}
	if len(args) > 0 {
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
	optionJDs(jds, options, viper.GetBool("report"))

	//output
	if viper.GetBool("output") {
		output(jds)
		os.Exit(0)
	}
	//reporter
	var reporter *exec.Reporter
	//main options
	mainOptions := []routine.Opt{routine.Interrupts(routine.DefaultCancelInterruptors...)}
	if needReport {
		n := viper.GetInt("repeat-times") * viper.GetInt("concurrent")
		reporter = exec.NewReporter(n)
		prepare := routine.ExecutorFunc(func(ctx context.Context) error {
			routine.Go(ctx, reporter)
			return nil
		})
		cleanup := routine.ExecutorFunc(func(ctx context.Context) error {
			reporter.Stop()
			reporter.Finalize()
			return nil
		})
		mainOptions = append(mainOptions, routine.Prepare(prepare), routine.Cleanup(cleanup))
	}
	//context
	ctx := context.Background()
	if httpConnections > 0 {
		httpclient.DefaultMaxConnsPerHost = httpConnections
		httpclient.DefaultMaxIdleConnsPerHost = httpConnections
		exec.WithTransport(ctx, httpclient.New().GetTransport())
	}

	jobs := NewJOBs(jds, reporter)
	jobs.Sort()

	exitForErr(routine.Main(
		withVerbose(ctx),
		jobs,
		mainOptions...))
}
