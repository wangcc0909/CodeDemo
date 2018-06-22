package model

import "time"

type Vote struct {
	ID            uint       `gorm:"primary_key" json:"id"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdateAt      time.Time  `json:"updateAt"`
	DeleteAt      *time.Time `sql:"index" json:"deleteAt"`
	EndAt         time.Time  `json:"endAt"`
	Name          string     `json:"name"`
	BrowseCount   int        `json:"browseCount"`
	CommentCount  int        `json:"commentCount"`
	CollectCount  int        `json:"collectCount"`
	Status        int        `json:"status"`
	Content       string     `json:"content"`
	HTMLContent   string     `json:"htmlContent"`
	ContentType   int        `json:"contentType"`
	Comments      []Comment  `gorm:"ForeignKey:SourceID" json:"comments"`
	UserID        uint       `json:"userId"`
	User          User       `json:"user"`
	LastUserID    uint       `json:"lastUserId"`
	LastUser      User       `json:"lastUser"`
	LastCommentAt *time.Time `json:"lastCommentAt"`
	VoteItems     []VoteItem `json:"voteItems"`
}

type VoteItem struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdateAt  time.Time  `json:"updateAt"`
	DeleteAt  *time.Time `sql:"index" json:"deleteAt"`
	Name      string     `json:"name"`
	Count     int        `json:"count"`
	VoteID    uint       `json:"voteId"`
}

//用户对那个投票项进行了投票
type UserVote struct {
	ID         uint       `gorm:"primary_key" json:"id"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdateAt   time.Time  `json:"updateAt"`
	DeleteAt   *time.Time `sql:"index" json:"deleteAt"`
	User       User       `json:"user"`
	UserID     uint       `json:"userId"`
	Vote       Vote       `json:"vote"`
	VoteID     uint       `json:"voteId"`
	VoteItemID uint       `json:"voteItemId"`
}

const (
	//进行中
	VoteUnderway = 1

	//已结束
	VoteOver = 2
)
