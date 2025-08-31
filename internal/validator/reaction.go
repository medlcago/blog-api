package validator

import (
	"blog-api/internal/models"
	"slices"

	"github.com/go-playground/validator/v10"
)

const (
	ReactionTag = "reaction"
)

func RegisterReactionValidation(v *validator.Validate) error {
	return v.RegisterValidation(ReactionTag, func(fl validator.FieldLevel) bool {
		reaction := fl.Field().String()
		return slices.Contains(models.AllowedReactions, models.ReactionType(reaction))
	})
}
