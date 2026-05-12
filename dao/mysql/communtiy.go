package mysql

import (
	"bluebell/models"
	"database/sql"

	"go.uber.org/zap"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	sqlstr := "select community_id,community_name from community"
	if err := db.Select(&communityList, sqlstr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("this is no community in db")
			err = nil
		}
	}
	return
}

// GetCommunityDetailByID 根据id查询社区分类详情
func GetCommunityDetailByID(id int64) (community *models.CommunityDetail, err error) {
	community = new(models.CommunityDetail)
	sqlstr := `select community_id,community_name,introduction,create_time 
			from community 
			where community_id=?
	`
	if err := db.Get(community, sqlstr, id); err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
		}
	}
	return community, err
}
