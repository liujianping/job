module github.com/x-mod/routine

require (
	github.com/cloudflare/cfssl v0.0.0-20190328212615-ea569c5aa1be
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/gorhill/cronexpr v0.0.0-20180427100037-88b0669f7d75
	github.com/rakyll/hey v0.1.1 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/x-mod/cronexpr v1.0.0
	github.com/x-mod/errors v0.1.0
)

replace github.com/x-mod/errors v0.1.0 => ../errors
