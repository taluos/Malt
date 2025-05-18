package errors

import (
	"google.golang.org/grpc/status"

	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

// ToGRPCError changes the error to a GRPC error
// if err is null return nil
// if err is *withCode, return the error with the code
// otherwise return the error with the unknown code
func ToGRPCError(e error) error {
	if e == nil {
		return e
	}
	var perr *withCode
	if As(e, &perr) {
		err := grpcStatus.Error(grpcCodes.Code(perr.code), perr.err.Error())
		return err
	}
	return grpcStatus.Error(grpcCodes.Unknown, e.Error())
}

// FromGRPCError changes the GRPC error to a error
// if err is null return nil
func FromGRPCError(e error) error {
	if e == nil {
		return e
	}
	st, ok := status.FromError(e)
	if !ok {
		return WithCode(100002, "")
	}

	return &withCode{
		err:  st.Err(),
		code: int(st.Code()),
	}
}
