package constants

// Response messages
const (
	MSGOK                  = "Success"
	MSGCreated             = "Created New Resource"
	MSGInternalServerError = "Internal Server Error"
	MSGNotFound            = "Resource Not Found"
	MSGUnauthenticated     = "Unauthenticated"
	MSGUnauthorized        = "Unauthorized"
	MSGBadRequest          = "Bad Request"
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
