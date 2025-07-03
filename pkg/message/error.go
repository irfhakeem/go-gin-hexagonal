package message

// Message
const (
	FAILED_INTERNAL_SERVER_ERROR    = "Internal server error"
	FAILED_UNAUTHORIZED             = "Unauthorized"
	FAILED_FORBIDDEN                = "Forbidden"
	FAILED_TOKEN_INVALID            = "Token invalid"
	FAILED_TOKEN_NOT_FOUND          = "Token not found"
	FAILED_TOKEN_EXPIRED            = "Token expired"
	FAILED_REGISTER_USER            = "Failed to register user"
	FAILED_LOGIN_USER               = "Failed to login user"
	FAILED_VERIFY_USER              = "Failed to verify user"
	FAILED_GET_AUTHORIZATION_HEADER = "Authorization header not found"
	FAILED_INVALID_REQUEST_FORMAT   = "Invalid request format"
	FAILED_INVALID_ID_FORMAT        = "Invalid ID format"
	FAILED_PASSWORD_INCORRECT       = "Current password is incorrect"

	FAILED_GET_ALL_USERS       = "Failed to get all users"
	FAILED_GET_USER_BY_ID      = "Failed to get user by id"
	FAILED_CREATE_USER         = "Failed to create user"
	FAILED_UPDATE_USER         = "Failed to update user"
	FAILED_DELETE_USER         = "Failed to delete user"
	FAILED_USER_ALREADY_EXISTS = "User already exists"
)
