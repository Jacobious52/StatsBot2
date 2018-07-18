package storage

import (
	"time"
)

type MessageKey string
type UserKey string
type ChatKey int64

type Messages map[time.Time]MessageKey
type Chat map[UserKey]Messages
type Model map[ChatKey]Chat

type MessageInfo struct {
	Chat      ChatKey
	Sender    UserKey
	Timestamp time.Time
}
