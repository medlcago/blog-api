package users

type UserResponse struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Deleted  bool   `json:"deleted"`
}
