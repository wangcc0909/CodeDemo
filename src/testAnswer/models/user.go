package models

import "time"

type User struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
	Friends   []User     `gorm:"many2many:user_friend;ForeignKey:ID;AssociationForeignKey:ID" json:"friends"`
}
