package backendmanager

import "gorm.io/gorm"

type QueueManager struct {
	DB *gorm.DB
}

type Queue struct {
	gorm.Model
	Username   string `json:"username"`
	StartTime  int64  `json:"start_time"`
	EndTime    int64  `json:"end_time"`
	UseProxies bool   `json:"use_proxies"`
}

func NewQueueManager() *QueueManager {
	return &QueueManager{}
}
