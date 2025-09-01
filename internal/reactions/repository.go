package reactions

import (
	"blog-api/internal/models"

	"gorm.io/gorm"
)

func GetReactionsAggregate(db *gorm.DB, targetType string, targetIDs []uint) (map[uint][]models.ReactionStat, error) {
	var results []models.ReactionStat

	err := db.Table("reactions r").
		Select("r.target_id as target_id, rt.name as type, rt.icon as icon, COUNT(*) as count").
		Joins("JOIN reaction_types rt on r.reaction_type_id = rt.id").
		Where("r.target_type = ? AND r.target_id IN ?", targetType, targetIDs).
		Group("r.target_id, rt.name, rt.icon").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	aggMap := make(map[uint][]models.ReactionStat)
	for _, r := range results {
		aggMap[r.TargetID] = append(aggMap[r.TargetID], models.ReactionStat{
			TargetID: r.TargetID,
			Type:     r.Type,
			Count:    r.Count,
			Icon:     r.Icon,
		})
	}
	return aggMap, nil
}

func GetUserReactions(db *gorm.DB, targetType string, targetIDs []uint, userID uint) (map[uint]*models.UserReaction, error) {
	var results []models.UserReaction

	err := db.Table("reactions r").
		Select(`r.target_id as target_id, rt.name as type, rt.icon as icon`).
		Joins("JOIN reaction_types rt ON r.reaction_type_id = rt.id").
		Where("r.target_type = ? AND r.target_id IN ? AND r.user_id = ?", targetType, targetIDs, userID).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	userMap := make(map[uint]*models.UserReaction)
	for _, r := range results {
		userMap[r.TargetID] = &r
	}

	return userMap, nil
}
