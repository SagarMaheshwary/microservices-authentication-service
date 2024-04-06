package constants

// gRPC response messages
const (
	OK                    = "Success"
	CREATED               = "Created New Resource"
	INTERNAL_SERVER_ERROR = "Internal Server Error"
	NOT_FOUND             = "Resource Not Found"
	UNAUTHENTICATED       = "Unauthenticated"
	UNAUTHORIZED          = "Unauthorized"
	BAD_REQUEST           = "Bad Request"
)

// gRPC metadata headers
const (
	HDR_AUTHORIZATION = "authorization"
)

const HDR_BEARER_PREFIX = "Bearer "

// Redis key prefix
const (
	RDS_TOKEN_BLACKLIST = "token-blacklist"
)
