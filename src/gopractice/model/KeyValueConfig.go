package model

import "time"

type KeyValueConfig struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdateAt  time.Time  `json:"updateAt"`
	DeleteAt  *time.Time `sql:"index" json:"deleteAt"`
	KeyName   string     `json:"key"`
	Value     string     `json:"value"`
}
