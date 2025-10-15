package utils

const (
	ValidationError = "VALIDATION_ERROR"

	Unauthorized       = "UNAUTHORIZED"
	InvalidCredentials = "INVALID_CREDENTIALS"
	InvalidToken       = "INVALID_TOKEN"
	ForbiddenAction    = "FORBIDDEN_ACTION"
	AccountLocked      = "ACCOUNT_LOCKED"
	TooManyRequests    = "TOO_MANY_REQUESTS"

	NotFound = "NOT_FOUND"
	Conflict = "CONFLICT"

	InternalError      = "INTERNAL_ERROR"
	ServiceUnavailable = "SERVICE_UNAVAILABLE"
	Timeout            = "TIMEOUT"
)
