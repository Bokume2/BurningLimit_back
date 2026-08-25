package schema

import "gorm.io/gorm"

type User struct {
	ID string `gorm:"primaryKey,check:id != ''"`
	gorm.Model
	Username    string `gorm:"uniqueIndex,check:username != ''"`
	Email       string `gorm:"uniqueIndex,check:email != ''"`
	Displayname string
	Description string
	Groups      []*Group `gorm:"many2many:user_groups;"`
	Roles       []*Role  `gorm:"many2many:user_roles;"`
	Tasks       []*Task  `gorm:"many2many:user_tasks;"`
}
