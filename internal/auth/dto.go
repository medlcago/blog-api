package auth

type RegisterUserInput struct {
	Username string `json:"username" validate:"username"`
	Password string `json:"password" validate:"required,min=6,max=60"`
}

type LoginUserInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TokenResponse struct {
	AccessToken          string `json:"access_token"`
	AccessTokenExpiresIn int    `json:"access_token_expires_in"`

	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
}
