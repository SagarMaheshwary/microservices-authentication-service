package constant

// Response messages
const (
	MessageOK                  = "Success"
	MessageCreated             = "Created New Resource"
	MessageBadRequest          = "Bad Request"
	MessageUnauthorized        = "Unauthorized"
	MessageForbidden           = "Forbidden"
	MessageNotFound            = "Resource Not Found"
	MessageInternalServerError = "Internal Server Error"
)

// gRPC metadata headers
const (
	HeaderAuthorization = "authorization"
)

const HeaderBearerPrefix = "Bearer "

// Redis key prefix
const (
	RedisTokenBlacklist = "token-blacklist"
)

const ServiceName = "Authentication Service"

const ExitFailure = 1
