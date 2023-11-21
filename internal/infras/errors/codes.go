package errors

type ErrorCode int

const (
	// OK Returned by the Code function on a nil error. It is not a valid
	// code for an error.
	OK ErrorCode = 0

	// Unknown The error could not be categorized.
	Unknown ErrorCode = 1

	// NotFound The resource was not found.
	NotFound ErrorCode = 2

	// AlreadyExists The resource exists, but it should not.
	AlreadyExists ErrorCode = 3

	// InvalidArgument A value given to a Go CDK API is incorrect.
	InvalidArgument ErrorCode = 4

	// Internal Something unexpected happened. Internal errors always indicate
	// bugs in the Go CDK (or possibly the underlying service).
	Internal ErrorCode = 5

	// Unimplemented The feature is not implemented.
	Unimplemented ErrorCode = 6

	// FailedPrecondition The system was in the wrong state.
	FailedPrecondition ErrorCode = 7

	// PermissionDenied The caller does not have permission to execute the specified operation.
	PermissionDenied ErrorCode = 8

	// ResourceExhausted Some resource has been exhausted, typically because a service resource limit
	// has been reached.
	ResourceExhausted ErrorCode = 9

	// Canceled The operation was canceled.
	Canceled ErrorCode = 10

	// DeadlineExceeded The operation timed out.
	DeadlineExceeded ErrorCode = 11

	// Unauthenticated The authentication operation failed.
	Unauthenticated ErrorCode = 12

	Unavailable ErrorCode = 13
)

func (i ErrorCode) String() string {
	switch i {
	case NotFound:
		return "NotFound"
	case AlreadyExists:
		return "AlreadyExists"
	case InvalidArgument:
		return "InvalidArgument"
	case Internal:
		return "Internal"
	case Unimplemented:
		return "Unimplemented"
	case FailedPrecondition:
		return "FailedPrecondition"
	case PermissionDenied:
		return "PermissionDenied"
	case ResourceExhausted:
		return "ResourceExhausted"
	case Canceled:
		return "Canceled"
	case DeadlineExceeded:
		return "DeadlineExceeded"
	case Unauthenticated:
		return "Unauthenticated"
	case Unavailable:
		return "Unavailable"
	}

	return "Unknown"
}
