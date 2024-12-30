package cache

import "context"

type FriendRelationCache interface {
	BatchDeleter
	DelFollowerUserIDs(ctx context.Context, userID string) FriendRelationCache
	DelSubscriberUserIDs(ctx context.Context, userID string) FriendRelationCache
	CloneFriendRelationCache() FriendRelationCache
	GetFollowerUserIDs(ctx context.Context, userID string) (userIDs []string, err error)
	GetSubscriberUserIDs(ctx context.Context, userID string) (userIDs []string, err error)
}
