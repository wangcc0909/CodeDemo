package model

import "time"

type Category struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdateAt  time.Time  `json:"updateAt"`
	DeleteAt  *time.Time `sql:"index" json:"deleteAt"`
	Name      string     `json:"name"`
	Sequence  int        `json:"sequence"`
	ParentID  int        `json:"parentId"`
}
