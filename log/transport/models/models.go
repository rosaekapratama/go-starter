package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"time"
)

type TransportLog struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	TraceID      string    `gorm:"type:varchar(32);not null"`
	SpanID       string    `gorm:"type:varchar(16);not null"`
	Type         string    `gorm:"type:varchar(16);not null"`
	Log          datatypes.JSON
	ErrorMessage *string   `gorm:"type:text"`
	ProcessDT    time.Time `gorm:"type:timestamptz;not null;default:now()"`
	ProcessBy    string    `gorm:"type:varchar(255);not null"`
}

type TransportRestLog struct {
	IsServer   bool    `json:"isServer"`
	IsRequest  bool    `json:"isRequest"`
	URL        string  `json:"url"`
	Method     string  `json:"method"`
	Headers    *string `json:"headers,omitempty"`
	Body       *string `json:"body,omitempty"`
	StatusCode *string `json:"statusCode,omitempty"`
}

type TransportGooglePubSubLog struct {
	IsPublisher bool    `json:"isPublisher"`
	TopicId     string  `json:"topicId"`
	MessageId   *string `json:"messageId,omitempty"`
	MessageData *string `json:"messageData,omitempty"`
}

type TransportSoapLog struct {
	IsServer   bool    `json:"isServer"`
	IsRequest  bool    `json:"isRequest"`
	URL        string  `json:"url"`
	SOAPAction string  `json:"SOAPAction"`
	Method     string  `json:"method"`
	Headers    string  `json:"headers"`
	Body       *string `json:"body,omitempty"`
	StatusCode *string `json:"statusCode,omitempty"`
}
