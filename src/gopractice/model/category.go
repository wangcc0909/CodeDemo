package model

import "time"

type Category struct {
	ID       uint       `gorm:"primary_key" json:"id"`
	CreateAt time.Time  `json:"createAt"`
	UpdateAt time.Time  `json:"updateAt"`
	DeleteAt *time.Time `sql:"index" json:"deleteAt"`
	Name     string     `json:"name"`
	Sequence int        `json:"sequence"`
	ParentID int        `json:"parentId"`
}
