package schema

import "gorm.io/gorm"

type TaskTag struct {
	gorm.Model
	Name    string  `gorm:"index:idx_name_in_group,unique"`
	GroupID uint    `gorm:"index:idx_name_in_group,unique"`
	Tasks   []*Task `gorm:"many2many:task_task_tags;"`
}
