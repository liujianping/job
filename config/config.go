package config

import (
	"strings"
	"time"
)

// Job:
//   name: "demo"
//   command:
//     name: "echo"
//     args:
//       - "hello"
//     envs:
//       - key: "key"
//         val: "val"
//   crontab: ""
//   repeat:
//     times: 0
//     interval: 100ms
//   retry: 3
//   timeout:
//     cmd: 3s
//     job: 1h

//JD struct config for job
type JD struct {
	Name       string        `yaml:"name"`
	Command    *Command      `yaml:"command"`
	HTTP       *HTTP         `yaml:"http"`
	Crontab    string        `yaml:"crontab"`
	Repeat     Repeat        `yaml:"repeat"`
	Concurrent int           `yaml:"concurrent"`
	Timeout    time.Duration `yaml:"timeout"`
	Guarantee  bool          `yaml:"guarantee"`
	Report     bool          `yaml:"report"`
	Order      Order         `yaml:"order"`
}

//Order struct
type Order struct {
	Precondition []string `yaml:"precondition"`
	Weight       int      `yaml:"weight"`
	Wait         bool     `yaml:"wait"`
}

//Command struct
type Command struct {
	Name    string        `yaml:"name"`
	Args    []string      `yaml:"args"`
	Envs    []KV          `yaml:"envs"`
	Retry   int           `yaml:"retry"`
	Timeout time.Duration `yaml:"timeout"`
}

//HTTP struct
type HTTP struct {
	Request  Request       `yaml:"request"`
	Response *Response     `yaml:"response"`
	Retry    int           `yaml:"retry"`
	Timeout  time.Duration `yaml:"timeout"`
}

//Request struct
type Request struct {
	Method  string            `yaml:"method"`
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers"`
	Query   map[string]string `yaml:"query"`
	Body    *Body             `yaml:"body"`
}

//Response struct
type Response struct {
	Status int   `yaml:"status"`
	Body   *Body `yaml:"body"`
}

//JSON struct
type JSON map[string]interface{}

//XML struct
type XML map[string]interface{}

//Text struct
type Text string

//Body struct
type Body struct {
	Text *Text `yaml:"text"`
	JSON *JSON `yaml:"json"`
	XML  *XML  `yaml:"xml"`
}

//KV struct
type KV struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

//Repeat struct
type Repeat struct {
	Times    int           `yaml:"times"`
	Interval time.Duration `yaml:"interval"`
}

func CommandJD() *JD {
	return &JD{
		Command:    &Command{},
		Repeat:     Repeat{Times: 1},
		Concurrent: 1,
		Report:     true,
	}
}

func HttpJD() *JD {
	return &JD{
		HTTP:       &HTTP{},
		Repeat:     Repeat{Times: 1},
		Concurrent: 1,
		Report:     true,
	}
}

func (jd *JD) String() string {
	if jd.Name != "" {
		return jd.Name
	}
	if jd.Command != nil {
		cmds := []string{jd.Command.Name}
		cmds = append(cmds, jd.Command.Args...)
		return strings.Join(cmds, " ")
	}
	if jd.HTTP != nil {
		if jd.HTTP.Request.URL == "" {
			return "http"
		}
		return jd.HTTP.Request.URL
	}
	return ""
}

//Option for JD
type Option func(*JD)

//Name opt
func Name(n string) Option {
	return func(jd *JD) {
		if len(n) > 0 {
			jd.Name = n
		}
	}
}

func CommandName(cmd string) Option {
	return func(jd *JD) {
		if len(cmd) > 0 {
			jd.Command.Name = cmd
		}
	}
}

func CommandArgs(args ...string) Option {
	return func(jd *JD) {
		if len(args) > 0 {
			jd.Command.Args = args
		}
	}
}

func CommandEnv(key, val string) Option {
	return func(jd *JD) {
		if len(key) > 0 {
			jd.Command.Envs = append(jd.Command.Envs, KV{Name: key, Value: val})
		}
	}
}

func CommandRetry(n int) Option {
	return func(jd *JD) {
		if n > 0 {
			jd.Command.Retry = n
		}
	}
}

func CommandTimeout(d time.Duration) Option {
	return func(jd *JD) {
		if d > 0 {
			jd.Command.Timeout = d
		}
	}
}

func Crontab(plan string) Option {
	return func(jd *JD) {
		if len(plan) > 0 {
			jd.Crontab = plan
		}
	}
}

func RepeatTimes(n int) Option {
	return func(jd *JD) {
		if n != 1 {
			jd.Repeat.Times = n
		}
	}
}

func RepeatInterval(d time.Duration) Option {
	return func(jd *JD) {
		if d > 0 {
			jd.Repeat.Interval = d
		}
	}
}

func Timeout(d time.Duration) Option {
	return func(jd *JD) {
		if d > 0 {
			jd.Timeout = d
		}
	}
}

func Concurrent(c int) Option {
	return func(jd *JD) {
		if c > 0 {
			jd.Concurrent = c
		}
	}
}

func Guarantee(g bool) Option {
	return func(jd *JD) {
		jd.Guarantee = g
	}
}

//JDs type
type JDs []*JD

func (js JDs) Len() int {
	return len(js)
}

func (js JDs) Less(i, j int) bool {
	return js[i].Order.Weight < js[j].Order.Weight
}

func (js JDs) Swap(i, j int) {
	js[i], js[j] = js[j], js[i]
}
