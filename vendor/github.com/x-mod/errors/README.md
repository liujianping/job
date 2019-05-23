errors
===
[![Build Status](https://travis-ci.org/x-mod/errors.svg?branch=master)](https://travis-ci.org/x-mod/errors) [![Go Report Card](https://goreportcard.com/badge/github.com/x-mod/errors)](https://goreportcard.com/report/github.com/x-mod/errors) [![Coverage Status](https://coveralls.io/repos/github/x-mod/errors/badge.svg?branch=master)](https://coveralls.io/github/x-mod/errors?branch=master) [![GoDoc](https://godoc.org/github.com/x-mod/errors?status.svg)](https://godoc.org/github.com/x-mod/errors) 

extension of errors for the following features:

- annotation error
- error with code, support grpc error code
- error by value
- error with stack

## Quick Start

````go

import "github.com/x-mod/errors"

//annotation error
e1 := errors.Annotate(err, "annotations")
e11 := errors.Annotatef(err, "annotations %s", "format")

//code error
e2 := errors.WithCode(err, YourCode)
//get error's code value
v2 := errors.ValueFrom(e2)

e3 := errors.WithStack(err)
//print stack info
fmt.Printf("%v", e3)

//pure value error
e4 := errors.ValueErr(5)
v4 := errors.ValueFrom(e4)

````
