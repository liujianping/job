package main

import (
	"github.com/liujianping/job/job"
	"github.com/x-mod/build"
	"github.com/x-mod/cmd"
)

func main() {
	c := cmd.Add(
		cmd.Name("job"),
		cmd.Main(job.Main),
	)
	c.Command.Use = "job [flags] [command args ...]"
	c.Command.Short = "Job, make your short-term command as a long-term job"
	c.Command.Example = `
	(simple)      $: job echo hello	
	(schedule)    $: job -s "* * * * *" -- echo hello
	(retry)       $: job -r 3 -- echox hello
	(repeat)      $: job -n 10 -i 100ms -- echo hello
	(concurrent)  $: job -c 10 -n 10 -- echo hello
	(timeout cmd) $: job -t 500ms -- sleep 1
	(timeout job) $: job -T 3s -r 4 -- sleep 1`

	c.Flags().DurationP("cmd-timeout", "t", 0, "job command timeout duration")
	c.Flags().DurationP("job-timeout", "T", 0, "job timeout duration")
	c.Flags().IntP("retry", "r", 0, "job command retry times when failed")
	c.Flags().IntP("concurrent", "c", 0, "job concurrent numbers ")
	c.Flags().IntP("repeat-times", "n", 1, "job repeat times, 0 means forever")
	c.Flags().DurationP("repeat-interval", "i", 0, "job repeat interval duration")
	c.Flags().StringP("schedule", "s", "", "job schedule in crontab format")

	cmd.Version(build.String())
	cmd.Execute()
}
