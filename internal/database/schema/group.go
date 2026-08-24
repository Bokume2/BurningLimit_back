package schema

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name        string
	Member      []*User `gorm:"many2many:user_groups;"`
	Description string
	Roles       []Role
	Tasks       []Task
}
