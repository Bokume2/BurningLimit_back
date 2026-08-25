package schema

import "gorm.io/gorm"

type TaskTag struct {
	gorm.Model
	Name    string
	GroupID uint
	Tasks   []*Task `gorm:"many2many:task_task_tags;"`
}
