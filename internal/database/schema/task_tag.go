package schema

import "gorm.io/gorm"

type TaskTag struct {
	gorm.Model
	Name  string
	Tasks []*Task `gorm:"many2many:task_task_tags;"`
}
