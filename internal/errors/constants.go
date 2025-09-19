package errors

var (
	ErrUsernameAlreadyExists = New(409, "username already exists")
	ErrInvalidCredentials    = New(401, "invalid credentials")
	ErrMissingToken          = New(401, "token is missing")
	ErrInvalidToken          = New(401, "invalid token")
	ErrUnauthorized          = New(401, "unauthorized")
	ErrNotFound              = New(404, "not found")
	ErrForbidden             = New(403, "forbidden")
	ErrInvalidFile           = New(400, "invalid file")
	ErrInvalidQuery          = New(400, "invalid query parameters")

	ErrIncorrectOldPassword = New(400, "incorrect old password")
	ErrNewPasswordSameAsOld = New(400, "new password cannot be the same as old password")

	ErrTwoFANotEnabled       = New(400, "2FA not enabled for user")
	ErrTwoFAAlreadyEnabled   = New(400, "2FA is already enabled")
	ErrInvalid2FACode        = New(400, "invalid 2FA code")
	ErrTwoFAFlowNotInitiated = New(400, "2FA login flow not initiated")
)
