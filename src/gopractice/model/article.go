package model

import "time"

type Article struct {
	ID            uint       `gorm:"primary_key" json:"id"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdateAt      time.Time  `json:"updateAt"`
	DeleteAt      *time.Time `sql:"index" json:"deleteAt"`
	Name          string     `json:"name"`
	BrowseCount   uint       `json:"browseCount"`
	CommentCount  uint       `json:"commentCount"`
	CollectCount  uint       `json:"collectCount"`
	Status        int        `json:"status"`
	Content       string     `json:"content"`
	HTMLContent   string     `json:"htmlContent"`
	ContentType   int        `json:"contentType"`
	Categories    []Category `gorm:"many2many:article_category;ForeignKey:ID;AssociationForeignKey:ID" json:"categories"`
	Comments      []Comment  `gorm:"ForeignKey:SourceID" json:"comments"`
	UserID        uint       `json:"userId"`
	User          User       `json:"user"`
	LastUserID    uint       `json:"lastUserId"`
	LastUser      User       `json:"lastUser"`
	LastCommentAt *time.Time `json:"lastCommentAt"`
}

const (
	//ArticleVerifying  文章正在审核
	ArticleVerifying = 1

	//文章审核成功
	ArticleVerifySuccess = 2

	//文章审核失败
	ArticleVerifyFail = 3
)

//最多能置顶的文章书
const MaxTopArticleCount = 4

type TopArticle struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdateAt  time.Time  `json:"updateAt"`
	DeleteAt  *time.Time `json:"deleteAt"`
	ArticleID uint       `json:"articleId"`
}
