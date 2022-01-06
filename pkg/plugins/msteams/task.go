package msteams

import (
	"gorm.io/gorm"
	"time"
)

const (
	PENDING string = "PENDING"
	SUCCESS string = "SUCCESS"
)

type Task struct {
	gorm.Model
	Status  string
	SendAt  time.Time
	Message string
	StartAt time.Time
}
