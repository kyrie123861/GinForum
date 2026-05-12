package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

func GetCommunity() ([]*models.Community, error) {
	// 查询数据库，查找到所有并返回
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(id)
}
