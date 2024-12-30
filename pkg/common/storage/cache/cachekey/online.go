package cachekey

import (
	"strings"
	"time"
)

const (
	OnlineKey           = "ONLINE:"
	LatestOnlineTimeKey = "ONLINETIME:"
	OnlineChannel       = "online_change"
	OnlineExpire        = time.Hour / 2
)

func GetOnlineKey(userID string) string {
	return OnlineKey + userID
}

func GetLatestOnlineTimeKey(userID string) string {
	return LatestOnlineTimeKey + userID
}

func GetOnlineKeyUserID(key string) string {
	return strings.TrimPrefix(key, OnlineKey)
}
