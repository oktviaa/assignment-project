package models

import "time"

type Score struct {
	ID         			uint   `gorm:"primaryKey"`
	AssignmentTitle     string `gorm:"not null;type:varchar(191)"`
	Description		    string `gorm:"not null;type:varchar(191)"`
	Score	    		int	   `gorm:"not null"`
	StudentID 			uint
	CreatedAt 			time.Time
	UpdatedAt 			time.Time
}