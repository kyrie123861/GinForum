package mysql

import (
	"bluebell/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	sqlstr := `insert into post(
	post_id,title,content,author_id,community_id)
	values(?,?,?,?,?)
	`
	_, err = db.Exec(sqlstr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

// GetPostById 根据ID查询单个帖子数据
func GetPostById(pid int64) (data *models.Post, err error) {
	post := new(models.Post)
	sqlstr := `select 
		post_id, title, content, author_id, community_id, create_time 
		from post 
		where post_id = ?
	`
	err = db.Get(post, sqlstr, pid)
	return post, err
}

// GetPostList 查询帖子列表数量
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlstr := `select 
		post_id, title, content, author_id, community_id, create_time 
		from post
		order by create_time 
		desc 
		limit ?, ?
	`
	posts = make([]*models.Post, 0, 2)
	err = db.Select(&posts, sqlstr, (page-1)*size, size)
	return
}

// GetPostByIDs 根据给定的id查询帖子数据
func GetPostByIDs(ids []string) (postList []*models.Post, err error) {
	sqlstr := `select post_id,title,content,author_id,community_id,create_time
		from post
		where post_id in (?)
		order by FIND_IN_SET(post_id,?)
		`
	query, args, err := sqlx.In(sqlstr, ids, strings.Join(ids, ","))
	if err != nil {
		return
	}

	query = db.Rebind(query)

	err = db.Select(&postList, query, args...)
	return

}

// deepseek修复redis 解决帖子投票问题出错问题
// GetAllPostsWithTime 获取所有帖子及其创建时间
func GetAllPostsWithTime() ([]*models.Post, error) {
	sqlStr := `SELECT post_id, create_time FROM post`
	rows, err := db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := new(models.Post)
		err := rows.Scan(
			&post.ID,
			&post.CreateTime,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// deepseek修复redis 解决帖子投票问题出错问题
// GetAllPosts 获取所有帖子
func GetAllPosts() ([]*models.Post, error) {
	sqlStr := `SELECT post_id, author_id, community_id, status, title, content, create_time, update_time FROM post`

	var posts []*models.Post
	err := db.Select(&posts, sqlStr)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
