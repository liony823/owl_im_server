package cachekey

const (
	FollowerUserIDsKey   = "FOLLOWER_USER_IDS:"
	SubscriberUserIDsKey = "SUBSCRIBER_USER_IDS:"
)

func GetFollowerUserIDsKey(ownerUserID string) string {
	return FollowerUserIDsKey + ownerUserID
}

func GetSubscriberUserIDsKey(ownerUserID string) string {
	return SubscriberUserIDsKey + ownerUserID
}
