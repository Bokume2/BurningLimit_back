package schema

import "gorm.io/gorm"

type RolePermissionPolicy struct {
	gorm.Model
	RoleID   uint
	Resource string
	Readable bool
	Writable bool
}
