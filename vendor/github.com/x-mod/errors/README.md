errors
===

extension of errors for the following features:

- annotation 
- code
- useful interfaces

## Quick Start

````go

import "github.com/x-mod/errors"

e1 := errors.Annotate(err, "annotations")

e2 := errors.WithCode(err, code1)

code2 := errors.CodeFrom(e2)

````
