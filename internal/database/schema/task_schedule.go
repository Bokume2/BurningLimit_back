package schema

import (
	"time"

	"gorm.io/gorm"
)

type TaskSchedule struct {
	gorm.Model
	TaskID      uint
	Start       time.Time
	End         time.Time `gorm:"check:end >= start"`
	Achievement uint      `gorm:"check:percentage >= 0 and percantage <= 100"`
}
