package model

import (
	"time"
)

type BookCategory struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdateAt  time.Time  `json:"updateAt"`
	DeleteAt  *time.Time `sql:"index" json:"deleteAt"`
	Name      string     `json:"name"`
	Sequence  int        `json:"sequence"`
	ParentID  uint       `json:"parentId"` //直接分父类的ID
}

type Book struct {
	ID             uint           `gorm:"primary_key" json:"id"`
	CreateAt       time.Time      `json:"createAt"`
	UpdateAt       time.Time      `json:"updateAt"`
	DeleteAt       *time.Time     `sql:"index" json:"deleteAt"`
	Name           string         `json:"name"`
	CoverURL       string         `json:"coverURL"`
	BrowseCount    uint           `json:"browseCount"`
	CommentCount   uint           `json:"commentCount"`
	CollectCount   uint           `json:"collectCount"`
	Status         string         `json:"status"`
	ReadLimits     string         `json:"readLimits"`
	Content        string         `json:"content"`
	HTMLContent    string         `json:"htmlContent"`
	ContentType    int            `json:"contentType"`
	Categories     []BookCategory `gorm:"many2many:book_category;ForeignKey:ID;AssociationForeignKey:ID" json:"categories"`
	Comments       []BookComment  `json:"comments"`
	UserID         uint           `json:"userId"`
	User           User           `json:"user"`
	Star           uint           `json:"star"`
	OneStarCount   uint           `json:"oneStarCount"`
	TwoStarCount   uint           `json:"twoStarCount"`
	ThreeStarCount uint           `json:"threeStarCount"`
	FourStarCount  uint           `json:"fourStarCount"`
	FiveStarCount  uint           `json:"fiveStarCount"`
	TotalStarCount uint           `json:"totalStarCount"`
}

type BookComment struct {
	ID          uint          `gorm:"primary_key" json:"id"`
	CreateAt    time.Time     `json:"createAt"`
	UpdateAt    time.Time     `json:"updateAt"`
	DeleteAt    *time.Time    `sql:"index" json:"deleteAt"`
	Status      string        `json:"status"`
	Star        uint          `json:"star"`
	Content     string        `json:"content"`
	HTMLContent string        `json:"htmlContent"`
	ContentType int           `json:"contentType"`
	ParentID    uint          `json:"parentId"`
	Parents     []BookComment `json:"parents"` //所有的父评论
	BookID      uint          `json:"bookId"`
	UserID      uint          `json:"userId"`
	User        User          `json:"user"`
}

type BookChapter struct {
	ID           uint          `gorm:"primary_key" json:"id"`
	CreateAt     time.Time     `json:"createAt"`
	UpdateAt     time.Time     `json:"updateAt"`
	DeleteAt     *time.Time    `sql:"index" json:"deleteAt"`
	Name         string        `json:"name"`
	BrowseCount  uint          `json:"browseCount"`
	CommentCount uint          `json:"commentCount"`
	Content      string        `json:"content"`
	HTMLContent  string        `json:"htmlContent"`
	ContentType  int           `json:"contentType"`
	Contents     []BookComment `json:"contents"`
	UserID       uint          `json:"userId"`
	User         User          `json:"user"`
	ParentID     uint          `json:"parentId"`
	BookID       uint          `json:"bookId"`
}

type BookChapterComment struct {
	ID          uint                 `gorm:"primary_key" json:"id"`
	CreateAt    time.Time            `json:"createAt"`
	UpdateAt    time.Time            `json:"updateAt"`
	DeleteAt    *time.Time           `sql:"index" json:"deleteAt"`
	Status      string               `json:"status"`
	Content     string               `json:"content"`
	HTMLContent string               `json:"htmlContent"`
	ContentType int                  `json:"contentType"`
	ParentID    uint                 `json:"parentId"`
	Parents     []BookChapterComment `json:"parents"` //所有的父评论
	BookID      uint                 `json:"bookId"`
	ChapterID   uint                 `json:"chapterId"`
	UserID      uint                 `json:"userId"`
	User        User                 `json:"user"`
}
