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
Usage:
  job [flags] [command args ...]

Examples:

(simple)      $: job echo hello
(schedule)    $: job -s "* * * * *" -- echo hello
(retry)       $: job -r 3 -- echox hello
(repeat)      $: job -n 10 -i 100ms -- echo hello
(concurrent)  $: job -c 10 -n 10 -- echo hello
(report)      $: job -R -- echo hello
(timeout cmd) $: job -t 500ms -R -- sleep 1
(timeout job) $: job -n 10 -i 500ms -T 3s -R -- echo hello
(job output)  $: job -n 10 -i 500ms -T 3s -o -- echo hello
(job config)  $: job -f /path/to/job.yaml

Flags:
  -e, --cmd-env stringToString     job command enviromental variables (default [])
  -r, --cmd-retry int              job command retry times when failed
  -d, --cmd-stdout-discard         job command stdout discard ?
  -t, --cmd-timeout duration       job command timeout duration
  -c, --concurrent int             job concurrent numbers
  -f, --config string              job config file path
  -G, --guarantee                  job guarantee mode enable ?
  -h, --help                       help for job
  -M, --metadata stringToString    job metadata definition (default [])
  -N, --name string                job name definition
  -o, --output                     job yaml config output enable ?
  -i, --repeat-interval duration   job repeat interval duration
  -w, --repeat-interval-nowait     job repeat interval nowait for current command done ?
  -n, --repeat-times int           job repeat times, 0 means forever (default 1)
  -R, --report                     job reporter enable ?
  -s, --schedule string            job schedule in crontab format
  -T, --timeout duration           job timeout duration
  -V, --verbose                    job verbose log enable ?
````

- ** Output Job ** 

````
$: job -n 10 -i 500ms -T 3s -o -- curl https://www.baidu.com
Job:
  name: ""
  command:
    shell:
      name: curl
      args:
      - https://www.baidu.com
      envs: []
    stdout: true
    retry: 0
    timeout: 0s
  guarantee: false
  crontab: ""
  repeat:
    times: 10
    interval: 500ms
  concurrent: 1
  timeout: 3s
  report: true
  order:
    precondition: []
    weight: 0
    wait: false
````

- ** Multple Job Config **

````yaml
Job:
  name: "echo"
  command:
    shell: 
      name: "echo"
      args: 
        - hello
        - job
      envs:
        - name: "key"
          value: "val"
    retry: 3
    timeout: 3s
    guarantee: false
  crontab: ""
  concurrent: 0
  repeat:
    times: 2
    interval: 100ms
  timeout: 10s
  report: true
  order:
    precondition: [""]
    weight: 4
    wait: false
---
Job:
  name: "http"
  command:
    retry: 3
    timeout: 3s
    stdout: true
    http:    
      request: 
        url: "https://www.baidu.com"
        method: GET
        # headers: 
        #   Content-Type: application/json
        # body:
        #   json:
        #     hello: "demo"
        #     person:
        #       name: jay
        #       hobby: football
  crontab: ""
  concurrent: 2
  repeat:
    times: 3
    interval: "10ms"
  timeout: 1h
  report: true
  order:
    weight: 3
    precondition: ["echo"]
    wait: false
````


````
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

- ** Local Report ** 

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

