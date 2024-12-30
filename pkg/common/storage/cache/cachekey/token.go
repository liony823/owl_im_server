package cachekey

import "github.com/openimsdk/open-im-server/v3/pkg/protocol/constant"

const (
	UidPidToken = "UID_PID_TOKEN_STATUS:"
)

func GetTokenKey(userID string, platformID int) string {
	return UidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
}
