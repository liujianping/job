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

````
$: job -c echo 
````