package reactions

import "gorm.io/gorm"

type ReactionAggregate struct {
	TargetID uint
	Type     string
	Count    int64
}

type UserReaction struct {
	TargetID uint
	Type     string
}

func GetReactionsAggregate(db *gorm.DB, targetType string, targetIDs []uint) (map[uint]map[string]int64, error) {
	var results []ReactionAggregate

	err := db.Table("reactions").
		Select("target_id, type, COUNT(*) as count").
		Where("target_type = ? AND target_id IN ?", targetType, targetIDs).
		Group("target_id, type").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// map[targetID] -> map[type]count
	aggMap := make(map[uint]map[string]int64)
	for _, r := range results {
		if _, ok := aggMap[r.TargetID]; !ok {
			aggMap[r.TargetID] = make(map[string]int64)
		}
		aggMap[r.TargetID][r.Type] = r.Count
	}

	return aggMap, nil
}

func GetUserReactions(db *gorm.DB, targetType string, targetIDs []uint, userID uint) (map[uint]string, error) {
	var results []UserReaction

	err := db.Table("reactions").
		Select("target_id, type").
		Where("target_type = ? AND target_id IN ? AND user_id = ?", targetType, targetIDs, userID).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	userMap := make(map[uint]string)
	for _, r := range results {
		userMap[r.TargetID] = r.Type
	}

	return userMap, nil
}
