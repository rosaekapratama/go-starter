package models

import (
	"github.com/google/uuid"
	"time"
)

type TransportRestLog struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	TraceID      string    `gorm:"type:varchar(32);not null"`
	SpanID       string    `gorm:"type:varchar(16);not null"`
	IsServer     bool      `gorm:"type:boolean;not null"`
	IsRequest    bool      `gorm:"type:boolean;not null"`
	URL          string    `gorm:"type:varchar(255);not null"`
	Method       string    `gorm:"type:varchar(20);not null"`
	Headers      *string   `gorm:"type:text"`
	Body         *string   `gorm:"type:text"`
	StatusCode   *string   `gorm:"type:varchar(3)"`
	ErrorMessage *string   `gorm:"type:text"`
	ProcessDT    time.Time `gorm:"type:timestamptz;not null;default:now()"`
	ProcessBy    string    `gorm:"type:varchar(255);not null"`
}
