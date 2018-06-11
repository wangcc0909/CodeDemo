package model

import "time"

type data struct {
	Title          string `json:"title"`
	CommentContent string `json:"commentContent"`
}

type Message struct {
	ID         uint       `gorm:"primary_key" json:"id"`
	CreateAt   time.Time  `json:"createAt"`
	UpdateAt   time.Time  `json:"updateAt"`
	DeleteAt   *time.Time `sql:"index" json:"deleteAt"`
	Type       string     `json:"type"`
	Readed     bool       `json:"readed"`
	FromUserId uint       `json:"fromUserId"`
	ToUserId   uint       `json:"toUserId"`
	FromUser   User       `json:"fromUser"`
	SourceId   uint       `json:"sourceId"`
	SourceName string     `json:"sourceName"`
	CommentId  uint       `json:"commentId"`
	Data       data       `json:"data"`
}

const (
	//MessageTypeCommentArticle 回复了话题
	MessageTypeCommentArticle = "messageTypeCommentArticle"

	//MessageTypeCommentVote  回复了投票
	MessageTypeCommentVote = "messageTypeCommentVote"

	//MessageTypeCommentComment  对回复进行了回复
	MessageTypeCommentComment = "messageTypeCommentComment"
)
