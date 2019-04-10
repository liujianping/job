package errors

import (
	"fmt"
)

type annotator struct {
	err        error
	annotation string
}

//Annotate an error with annotation
func Annotate(err error, annotation string) error {
	if err != nil {
		return &annotator{err: err, annotation: annotation}
	}
	return nil
}

//Annotatef an error with annotation
func Annotatef(err error, format string, args ...interface{}) error {
	if err != nil {
		return &annotator{err: err, annotation: fmt.Sprintf(format, args...)}
	}
	return nil
}

//Error annotator implemention
func (err *annotator) Error() string {
	return fmt.Sprintf("%s: %s", err.annotation, err.err.Error())
}

//Cause annotator implemention
func (e *annotator) Cause() error {
	return e.err
}
