package schema

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username    string
	Email       string
	Displayname string
	Description string
	Groups      []*Group `gorm:"many2many:user_groups;"`
	Roles       []*Role  `gorm:"many2many:user_roles;"`
	Tasks       []*Task  `gorm:"many2many:user_tasks;"`
}
