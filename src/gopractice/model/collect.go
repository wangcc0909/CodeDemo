package model

import "time"

type Folder struct {
	ID       uint       `gorm:"primary_key" json:"id"`
	CreateAt time.Time  `json:"createAt"`
	UpdateAt time.Time  `json:"updateAt"`
	DeleteAt *time.Time `sql:"index" json:"deleteAt"`
	Name     string     `json:"name"`
	UserID   uint       `json:"userId"`
	ParentID uint       `json:"parentId"`
}

//收藏
type Collect struct {
	ID         uint       `gorm:"primary_key" json:"id"`
	CreateAt   time.Time  `json:"createAt"`
	UpdateAt   time.Time  `json:"updateAt"`
	DeleteAt   *time.Time `sql:"index" json:"deleteAt"`
	UserID     uint       `json:"userId"`
	SourceName string     `json:"sourceName"` //用来区分是对话题，还是对投票进行收藏
	SourceID   uint       `json:"sourceId"`
	FolderID   uint       `json:"folderId"`
	Folder     Folder     `json:"folder"`
}

const (
	//收藏文章
	CollectSourceArticle = "CollectSourceArticle"
	//收藏投票
	CollectSourceVote = "CollectSourceVote"
)

const MaxFolderCount = 20
