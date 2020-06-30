job
===
[![GoDoc](https://godoc.org/github.com/liujianping/job?status.svg)](https://godoc.org/github.com/liujianping/job) [![Go Report Card](https://goreportcard.com/badge/github.com/liujianping/job)](https://goreportcard.com/report/github.com/liujianping/job) [![Build Status](https://travis-ci.org/liujianping/job.svg?branch=master)](https://travis-ci.org/liujianping/job) [![Version](https://img.shields.io/github/tag/liujianping/job.svg)](https://github.com/liujianping/job/releases) [![Coverage Status](https://coveralls.io/repos/github/liujianping/job/badge.svg?branch=master)](https://coveralls.io/github/liujianping/job?branch=master)

make your short-term command as a long-term job

## Install

Brew install

````bash
$: brew tap liujianping/tap && brew install job
````

OR 

````bash
$: git clone https://github.com/liujianping/job.git
$: cd job 
$: go build -mod vendor
````

## Usage

````bash

$: job -h
Job, make your short-term command as a long-term job

Usage:
  job [flags] [command args ...]

Examples:

	(simple)      $: job echo hello
	(schedule)    $: job -s "* * * * *" -- echo hello
	(retry)       $: job -r 3 -- echox hello
	(repeat)      $: job -n 10 -i 100ms -- echo hello
	(concurrent)  $: job -c 10 -n 10 -- echo hello
	(timeout cmd) $: job -t 500ms -- sleep 1
	(timeout job) $: job -T 3s -r 4 -- sleep 1

Flags:
  -t, --cmd-timeout duration       job command timeout duration
  -c, --concurrent int             job concurrent numbers
  -h, --help                       help for job
  -T, --job-timeout duration       job timeout duration
  -i, --repeat-interval duration   job repeat interval duration
  -n, --repeat-times int           job repeat times, 0 means forever (default 1)
  -r, --retry int                  job command retry times when failed
  -s, --schedule string            job schedule in crontab format
      --version                    version for job
````
