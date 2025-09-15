package auth

type RegisterUserInput struct {
	Username string `json:"username" validate:"username"`
	Password string `json:"password" validate:"required,min=6,max=60"`
}

type LoginUserInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type TokenResponse struct {
	AccessToken          string `json:"access_token"`
	AccessTokenExpiresIn int    `json:"access_token_expires_in"`

	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
}

type LoginResponse struct {
	Token       *TokenResponse `json:"token,omitempty"`
	Requires2FA bool           `json:"requires_2fa"`
	Message     string         `json:"message,omitempty"`
}

type ChangePasswordInput struct {
	OldPassword string `form:"old_password" validate:"required"`
	NewPassword string `form:"new_password" validate:"required,min=6,max=60"`
}

type TwoFASetupResponse struct {
	QRCode  string `json:"qr_code"`
	Message string `json:"message,omitempty"`
}

type Verify2FAInput struct {
	Code string `json:"code" validate:"required"`
}

type Login2FAInput struct {
	Username string `json:"username" validate:"username"`
	Code     string `json:"code" validate:"required"`
}
