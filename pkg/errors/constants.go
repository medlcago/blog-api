package errors

var (
	ErrUsernameAlreadyExists = New(409, "username already exists")
	ErrInvalidCredentials    = New(401, "invalid credentials")
	ErrMissingToken          = New(401, "token is missing")
	ErrUnauthorized          = New(401, "unauthorized")
	ErrNotFound              = New(404, "not found")
	ErrInvalidQuery          = New(400, "invalid query parameters")
)
