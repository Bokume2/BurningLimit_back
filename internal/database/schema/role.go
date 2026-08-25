package schema

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name               string  `gorm:"index:idx_name_in_group,unique"`
	GroupID            uint    `gorm:"index:idx_name_in_group,unique"`
	Users              []*User `gorm:"many2many:user_roles;"`
	Tasks              []*Task `gorm:"many2many:role_tasks;"`
	PermissionLevel    uint
	PermissionPolicies []RolePermissionPolicy
}
