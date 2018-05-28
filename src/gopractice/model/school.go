package model

import "time"

type School struct {
	ID         uint       `gorm:"primary_key" json:"id"`
	CreateAt   time.Time  `json:"createAt"`
	UpdateAt   time.Time  `json:"updateAt"`
	DeleteAt   *time.Time `sql:"index" json:"deleteAt"`
	Name       string     `json:"name"`
	Speciality string     `json:"speciality"`
	UserID     uint       `json:"userId"`
}

const MaxSchoolNameLen = 200
const MaxSchoolSpecialityLen = 200
