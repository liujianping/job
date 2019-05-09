package config

import (
	"net/url"
	"strings"
	"time"
)

//JD struct config for job
type JD struct {
	Name       string                 `yaml:"name"`
	Metadata   map[string]interface{} `yaml:"metadata,omitempty"`
	Command    Command                `yaml:"command"`
	Guarantee  bool                   `yaml:"guarantee"`
	Crontab    string                 `yaml:"crontab"`
	Repeat     Repeat                 `yaml:"repeat"`
	Concurrent int                    `yaml:"concurrent"`
	Timeout    time.Duration          `yaml:"timeout"`
	Report     bool                   `yaml:"report"`
	Order      Order                  `yaml:"order"`
}

//Order struct
type Order struct {
	Precondition []string `yaml:"precondition"`
	Weight       int      `yaml:"weight"`
	Wait         bool     `yaml:"wait"`
}

//Command struct
type Command struct {
	Shell   *Shell        `yaml:"shell,omitempty"`
	HTTP    *HTTP         `yaml:"http,omitempty"`
	Stdout  bool          `yaml:"stdout"`
	Retry   int           `yaml:"retry"`
	Timeout time.Duration `yaml:"timeout"`
}

//Shell struct
type Shell struct {
	Name string   `yaml:"name"`
	Args []string `yaml:"args"`
	Envs []KV     `yaml:"envs"`
}

//HTTP struct
type HTTP struct {
	Request Request `yaml:"request"`
}

//Request struct
type Request struct {
	Method  string            `yaml:"method"`
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers"`
	Queries map[string]string `yaml:"queries"`
	Body    *Body             `yaml:"body"`
}

//JSON struct
type JSON map[string]interface{}

//XML struct
type XML map[string]interface{}

//Text struct
type Text string

//Body struct
type Body struct {
	Text *Text       `yaml:"text"`
	JSON *JSON       `yaml:"json"`
	XML  *XML        `yaml:"xml"`
	Form *url.Values `yaml:"form"`
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

//CommandJD default
func CommandJD() *JD {
	return &JD{
		Metadata: map[string]interface{}{},
		Command: Command{
			Shell: &Shell{},
		},
		Repeat:     Repeat{Times: 1},
		Concurrent: 1,
		Report:     false,
	}
}

//HTTPCommandJD default
func HTTPCommandJD() *JD {
	return &JD{
		Command: Command{
			HTTP: &HTTP{},
		},
		Repeat:     Repeat{Times: 1},
		Concurrent: 1,
		Report:     false,
	}
}

func (jd *JD) String() string {
	if jd.Name != "" {
		return jd.Name
	}
	if jd.Command.Shell != nil {
		cmds := []string{jd.Command.Shell.Name}
		cmds = append(cmds, jd.Command.Shell.Args...)
		return strings.Join(cmds, " ")
	}
	if jd.Command.HTTP != nil {
		if jd.Command.HTTP.Request.URL == "" {
			return "http"
		}
		return jd.Command.HTTP.Request.URL
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

//CommandName opt
func CommandName(cmd string) Option {
	return func(jd *JD) {
		if len(cmd) > 0 {
			jd.Command.Shell.Name = cmd
		}
	}
}

//CommandArgs opt
func CommandArgs(args ...string) Option {
	return func(jd *JD) {
		if len(args) > 0 {
			jd.Command.Shell.Args = append(jd.Command.Shell.Args, args...)
		}
	}
}

//CommandEnv opt
func CommandEnv(key, val string) Option {
	return func(jd *JD) {
		if len(key) > 0 {
			jd.Command.Shell.Envs = append(jd.Command.Shell.Envs, KV{Name: key, Value: val})
		}
	}
}

//CommandRetry opt
func CommandRetry(n int) Option {
	return func(jd *JD) {
		if n > 0 {
			jd.Command.Retry = n
		}
	}
}

//CommandTimeout opt
func CommandTimeout(d time.Duration) Option {
	return func(jd *JD) {
		if d > 0 {
			jd.Command.Timeout = d
		}
	}
}

//Guarantee opt
func Guarantee(g bool) Option {
	return func(jd *JD) {
		jd.Guarantee = g
	}
}

//CommandStdoutDiscard opt
func CommandStdoutDiscard(m bool) Option {
	return func(jd *JD) {
		jd.Command.Stdout = !m
	}
}

//Crontab opt
func Crontab(plan string) Option {
	return func(jd *JD) {
		if len(plan) > 0 {
			jd.Crontab = plan
		}
	}
}

//RepeatTimes opt
func RepeatTimes(n int) Option {
	return func(jd *JD) {
		if n != 1 {
			jd.Repeat.Times = n
		}
	}
}

//RepeatInterval opt
func RepeatInterval(d time.Duration) Option {
	return func(jd *JD) {
		if d > 0 {
			jd.Repeat.Interval = d
		}
	}
}

//Timeout opt
func Timeout(d time.Duration) Option {
	return func(jd *JD) {
		if d > 0 {
			jd.Timeout = d
		}
	}
}

//Concurrent opt
func Concurrent(c int) Option {
	return func(jd *JD) {
		if c != 0 {
			jd.Concurrent = c
		}
	}
}

//Metadata opt
func Metadata(key string, val interface{}) Option {
	return func(jd *JD) {
		if len(key) > 0 {
			jd.Metadata[key] = val
		}
	}
}
