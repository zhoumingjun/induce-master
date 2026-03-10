package service

import (
	"time"

	"github.com/google/uuid"
)

// GenerateUUID 生成 UUID
func GenerateUUID() string {
	return uuid.New().String()
}

func getCurrentTime() time.Time {
	return time.Now()
}
