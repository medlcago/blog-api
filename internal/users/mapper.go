package users

import "blog-api/internal/models"

func MapUserToResponse(user models.User) *UserResponse {
	return &UserResponse{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email.String,
		Deleted:  user.DeletedAt.Valid,
		Avatar:   user.Avatar.String,
	}
}
