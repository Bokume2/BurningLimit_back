package schema

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Name        string
	GroupID     uint
	Users       []*User `gorm:"many2many:user_tasks;"`
	Description string
	ParentID    *uint
	Children    []Task     `gorm:"foreignKey:ParentID"`
	Tags        []*TaskTag `gorm:"many2many:task_task_tags;"`
	Priority    string     `gorm:"type:enum('critical','high','middle','low');"`
	Roles       []Role     `gorm:"many2many:role_tasks;"`
	Schedule    []TaskSchedule
	DeadLine    *time.Time
}
