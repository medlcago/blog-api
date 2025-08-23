package posts

import (
	"blog-api/internal/models"
	goerrors "errors"
)

func ValidatePostEntities(post models.Post) error {
	contentLength := len(post.Content)
	entities := post.Entities

	for i, entity := range entities {
		if entity.Type == "link" && entity.URL == nil {
			return goerrors.New("type link: url is required")
		}
		if entity.Type != "link" && entity.URL != nil {
			return goerrors.New("url should only be provided for link type")
		}
		if entity.Offset+entity.Length > contentLength {
			return goerrors.New("entity range is out of bounds")
		}

		for j := i + 1; j < len(entities); j++ {
			otherEntity := entities[j]
			if entity.Offset == otherEntity.Offset && entity.Length == otherEntity.Length && entity.Type == otherEntity.Type {
				return goerrors.New("duplicate entity detected: type=" + entity.Type)
			}
		}
	}
	return nil
}
