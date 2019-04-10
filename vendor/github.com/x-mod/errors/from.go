package errors

import "google.golang.org/grpc/status"

//CodeFrom get code from the error
//support code from grpc status
func CodeFrom(err error) Code {
	if err != nil {
		for err != nil {
			//from grpc status
			if st, ok := status.FromError(err); ok {
				return &errorCode{value: int32(st.Code()), message: st.Message()}
			}
			//from error coder implement
			if cd, ok := err.(coder); ok {
				return &errorCode{value: int32(cd.Value()), message: err.Error()}
			}
			cause, ok := err.(causer)
			if !ok {
				break
			}
			err = cause.Cause()
		}
	}
	return nil
}

//ValueFrom get code value from error
//-1 means null code value
//0 means OK
func ValueFrom(err error) int {
	if err != nil {
		code := CodeFrom(err)
		if code != nil {
			return int(code.Value())
		}
		return -1
	}
	return 0
}

//CauseFrom get original error
func CauseFrom(err error) error {
	if err != nil {
		if cause, ok := err.(causer); ok {
			return cause.Cause()
		}
	}
	return err
}
