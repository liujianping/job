package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sort"

	"github.com/x-mod/httpclient"

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

	needReport := false
	httpShared := 0
	for _, jd := range jds {
		for _, opt := range options {
			opt(jd)
			if viper.GetBool("report") {
				jd.Report = true
			}
			if jd.Report {
				needReport = true
			}
			if jd.Command.HTTP != nil {
				httpShared = httpShared + jd.Concurrent
			}
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

	ctx := context.TODO()
	if viper.GetBool("verbose") {
		log := logrus.New()
		log.SetLevel(logrus.TraceLevel)
		ctx = routine.WithLogger(ctx, log)
	}
	if httpShared > 0 {
		ctx = exec.WithTransport(ctx, httpclient.NewHTTPTransport(httpclient.MaxIdleConnections(httpShared)))
	}
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := routine.Main(
		ctx,
		routine.ExecutorFunc(func(ctx context.Context) error {
			jobs := make(map[string]chan error, len(jds))
			var tail chan error
			for _, jd := range jds {
				if len(jd.Order.Precondition) == 0 {
					job := exec.NewJob(jd, reporter)
					ch := routine.Go(ctx, job)
					jobs[job.String()] = ch
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
					tail = ch
					if jd.Order.Wait {
						<-ch
					}
				}
			}
			return <-tail
		}),
		mainOptions...)
	if err != nil {
		fmt.Println("job failed:", err)
	}
	os.Exit(errors.ValueFrom(err))
}
