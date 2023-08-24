package models

import "time"

type Student struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"not null;type:varchar(191)"`
	Age	   	   int	  `gorm:"not null"`
	Scores 	   []Score
	CreatedAt time.Time
	UpdatedAt time.Time
}