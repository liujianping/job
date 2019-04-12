job
===

make your short-term command as a long-term job

## Install

````
$: go get -u github.com/liujianping/job
````

## Usage

````shell

$: job -h
Job, make your short-term command as a long-term job

Usage:
  job [flags]

Flags:
  -e, --cmd-env stringToString     job command enviromental variables (default [])
  -r, --cmd-retry int              job command retry times when failed
  -t, --cmd-timeout duration       job command timeout duration
  -C, --command string             job command path name
  -c, --concurrent int             job concurrent numbers  (default 1)
  -f, --config string              job config file path
  -s, --crontab string             job schedule plan in crontab format
  -G, --guarantee                  job guarantee mode enable ?
  -h, --help                       help for job
  -N, --name string                job name
  -i, --repeat-interval duration   job repeat interval duration
  -n, --repeat-times int           job repeat times, 0 means forever (default 1)
  -R, --report                     job reporter enable ?
  -T, --timeout duration           job timeout duration
````
## Examples

- **Crontab**

````
$: job -C echo hello -s "* * * * *" 

or

$: job -s "* * * * *" -- echo hello

````

- **Retry** when command failed

````

$: job -r 3 -- echox hello

````

- **Repeat** as you like 

````
$: job -n 10 -i 500ms -- echo hello

````

- **Concurrent**

````
$: job -n 10 -i 500ms -c 5 -- echo hello

````

- **Timeout** 

  - command timeout

````
$: job -t 500ms -- sleep 1
````

  - job timeout
  
````
$: job -n 0 -T 10s -- sleep 1
````

- **Yaml** config jobs in yaml format

````yaml
Job:
  name: "demo"
  command: 
    name: "echo"
    args: 
      - "hello"
      - "world"
    envs:
      - name: "key"
        value: "val"
    retry: 3
    timeout: 3s
  crontab: ""
  concurrent: 0
  repeat:
    times: 10
    interval: 100ms
  timeout: 1h
  guarantee: false
  report: true
  order:
    precondition: [""]
    weight: 4
    wait: false
---
Job:
  name: "work"
  http:
    retry: 3
    timeout: 3s
    request: 
      url: "http://liujianping.github.io"
      method: post
      headers: 
        Content-Type: application/json
        Authorization: Bearer {{env "AUTH_TOKEN"}}
      body:
        json:
          aa: "ddd"
          bb: false
    response:
      status: 200
      body:
        json:
          aa: "ddd"
          bb: false
  crontab: ""
  concurrent: 0
  repeat:
    times: 0
    interval: 100ms
  timeout: 1h
  guarantee: false
  report: false
  order:
    weight: 3
    precondition: [""]
    wait: false
````

run jobs:

````shell

$: job -f jobs.yaml 

````

- Report

````
$: job -n 10 -i 500ms -c 5 -R -- echo hello

Uptime:	5.1037 secs

Summary:
  Total:	5.1029 secs
  Slowest:	0.0091 secs
  Fastest:	0.0036 secs
  Average:	0.0068 secs
  Op/sec:	9.7983

  Total data:	210 bytes
  Size/Resp:	210 bytes

Response time histogram:
  0.004 [1]	|■■■■
  0.004 [0]	|
  0.005 [2]	|■■■■■■■■■
  0.005 [2]	|■■■■■■■■■
  0.006 [4]	|■■■■■■■■■■■■■■■■■■
  0.006 [9]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.007 [8]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.007 [7]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.008 [6]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.009 [6]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.009 [5]	|■■■■■■■■■■■■■■■■■■■■■■


Latency distribution:
  10% in 0.0056 secs
  25% in 0.0061 secs
  50% in 0.0068 secs
  75% in 0.0080 secs
  90% in 0.0086 secs
  95% in 0.0088 secs
  0% in 0.0000 secs

Code distribution:
  [0]	50 responses
````


- Job Timeout

````
$: job -C echo xxxx -n 0 -i 50ms -T 5s -R
Uptime:	5.0025 secs

Summary:
  Total:	5.0015 secs
  Slowest:	0.0051 secs
  Fastest:	0.0022 secs
  Average:	0.0041 secs
  Op/sec:	17.5946

  Total data:	528 bytes
  Size/Resp:	528 bytes

Response time histogram:
  0.002 [1]	|■
  0.003 [0]	|
  0.003 [0]	|
  0.003 [1]	|■
  0.003 [2]	|■
  0.004 [1]	|■
  0.004 [17]	|■■■■■■■■■■■■
  0.004 [56]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.005 [6]	|■■■■
  0.005 [2]	|■
  0.005 [2]	|■


Latency distribution:
  10% in 0.0039 secs
  25% in 0.0040 secs
  50% in 0.0041 secs
  75% in 0.0042 secs
  90% in 0.0043 secs
  95% in 0.0047 secs
  0% in 0.0000 secs

Code distribution:
  [0]	88 responses



exit status 255
````
