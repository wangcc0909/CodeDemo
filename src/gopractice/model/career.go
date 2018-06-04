package model

import "time"

type Career struct {
	ID       uint       `gorm:"primary_key" json:"id"`
	CreateAt time.Time  `json:"createAt"`
	UpdateAt time.Time  `sql:"index" json:"updateAt"`
	DeleteAt *time.Time `json:"deleteAt"`
	Company  string     `json:"company"`
	Title    string     `json:"title"` //职位
	UserID   uint       `json:"userId"`
}

const MaxCareerCompanyLen = 200

const MaxCareerTitleLen = 200
