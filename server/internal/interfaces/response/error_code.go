package response

type ErrorCode string

const (
	CodeInvalidPayload      ErrorCode = "INVALID_PAYLOAD"
	CodeDuplicateRoutine    ErrorCode = "DUPLICATE_ROUTINE"
	CodeInternalServerError ErrorCode = "INTERNAL_SERVER_ERROR"
	CodeNotFound            ErrorCode = "NOT_FOUND"
)
