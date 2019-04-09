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
  -a, --arg stringArray            job command argments
  -C, --cmd string                 job command path
  -t, --cmd-timeout duration       timeout duration for command
  -c, --concurrent int             job command concurrent numbers  (default 1)
      --config string              config file (default is $HOME/.job.yaml)
      --crontab string             job command schedule plan in crontab format
  -e, --env stringToString         job command enviromental variables (default [])
  -G, --guarantee                  job executing in guarantee mode
  -h, --help                       help for job
  -T, --job-timeout duration       timeout duration for the job
      --name string                job name
  -p, --payload bytesBase64        job command custom payload
  -P, --pipline                    job executing in pipeline mode
  -i, --repeat-interval duration   job command repeat interval duration
  -n, --repeat-times int           job command repeat times, 0 means forever (default 1)
  -R, --report                     job generate report
  -r, --retry-times int            job command retry times when failed
````

## Examples

- Crontab

````
$: job -C echo hello --crontab "* * * * *"
````

- Retry when command failed

````
$: job -C echox hello -r 3 
````

- Repeat

````
$: job -C echo hello -n 10 -i 500ms
````

- Concurrent

````
$: job -C echo hello -n 10 -i 500ms -c 5
````

- Command Timeout

````
$: job -C sleep -a 1 -t 500ms 
````

- Report

````
$: job -C echo hello -n 10 -i 500ms -c 5 -R 

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
