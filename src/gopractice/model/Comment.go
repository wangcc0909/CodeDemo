package model

import "time"

//评论
type Comment struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	CreateAt    time.Time  `json:"createAt"`
	UpdateAt    time.Time  `json:"updateAt"`
	DeleteAt    *time.Time `sql:"index" json:"deleteAt"`
	Status      int        `json:"status"`
	Content     string     `json:"content"`
	HTMLContent string     `json:"htmlContent"`
	ContentType int        `json:"contentType"`
	ParentID    uint       `json:"parentId"`   //直接父评论的Id
	Parents     []Comment  `json:"parents"`    //所有的父评论
	SourceName  string     `json:"sourceName"` //用来区分是对话还是对投票进行评论
	SourceID    uint       `json:"sourceId"`   //话题或者投票的ID
	UserID      uint       `json:"userId"`
	User        User       `json:"user"`
}

const (
	//对话题进行评论
	CommentSourceArticle = "article"

	//对投票进行评论
	CommentSourceVote = "vote"
)

const (
	//审核中
	CommentVertifying = 1

	//审核通过
	CommentVertifySuccess = 2

	//审核未通过
	CommentVertifyFail = 3
)

//评论的最大长度
const MaxCommentLen = 500
