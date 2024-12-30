package redis

import (
	"context"
	"time"

	"github.com/dtm-labs/rockscache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/tools/log"
	"github.com/redis/go-redis/v9"
)

const (
	followingRelatedUserIDsExpireTime  = time.Second * 60 * 60 * 12
	subscribedRelatedUserIDsExpireTime = time.Second * 60 * 60 * 12
)

type friendRelationCacheRedis struct {
	cache.BatchDeleter
	followingRelatedUserIDsExpireTime  time.Duration
	subscribedRelatedUserIDsExpireTime time.Duration
	rcClient                           *rockscache.Client
	friendRelationDatabase               database.FriendRelation
}

func NewFriendRelationCacheRedis(rdb redis.UniversalClient, localCache *config.LocalCache, friendRelationDatabase database.FriendRelation, options *rockscache.Options) cache.FriendRelationCache {
	batchHandler := NewBatchDeleterRedis(rdb, options, []string{localCache.Friend.Topic})
	b := localCache.Friend
	log.ZDebug(context.Background(), "user relation local cache init", "Topic", b.Topic, "SlotNum", b.SlotNum, "SlotSize", b.SlotSize, "enable", b.Enable())
	return &friendRelationCacheRedis{
		BatchDeleter:                       batchHandler,
		followingRelatedUserIDsExpireTime:  followingRelatedUserIDsExpireTime,
		subscribedRelatedUserIDsExpireTime: subscribedRelatedUserIDsExpireTime,
		rcClient:                           rockscache.NewClient(rdb, *options),
		friendRelationDatabase:               friendRelationDatabase,
	}
}

func (u *friendRelationCacheRedis) CloneFriendRelationCache() cache.FriendRelationCache {
	return &friendRelationCacheRedis{
		BatchDeleter:                       u.BatchDeleter.Clone(),
		followingRelatedUserIDsExpireTime:  u.followingRelatedUserIDsExpireTime,
		subscribedRelatedUserIDsExpireTime: u.subscribedRelatedUserIDsExpireTime,
		rcClient:                           u.rcClient,
		friendRelationDatabase:               u.friendRelationDatabase,
	}
}

func (u *friendRelationCacheRedis) getFollowerUserIDsKey(userID string) string {
	return cachekey.GetFollowerUserIDsKey(userID)
}

func (u *friendRelationCacheRedis) getSubscriberUserIDsKey(userID string) string {
	return cachekey.GetSubscriberUserIDsKey(userID)
}

func (u *friendRelationCacheRedis) GetFollowerUserIDs(ctx context.Context, userID string) (userIDs []string, err error) {
	return getCache(
		ctx,
		u.rcClient,
		u.getFollowerUserIDsKey(userID),
		u.followingRelatedUserIDsExpireTime,
		func(ctx context.Context) ([]string, error) {
			return u.friendRelationDatabase.GetFollowerUserIDs(ctx, userID)
		},
	)
}

func (u *friendRelationCacheRedis) GetSubscriberUserIDs(ctx context.Context, userID string) (userIDs []string, err error) {
	return getCache(
		ctx,
		u.rcClient,
		u.getSubscriberUserIDsKey(userID),
		u.subscribedRelatedUserIDsExpireTime,
		func(ctx context.Context) ([]string, error) {
			return u.friendRelationDatabase.GetSubscriberUserIDs(ctx, userID)
		},
	)
}

func (u *friendRelationCacheRedis) DelFollowerUserIDs(_ context.Context, userID string) cache.FriendRelationCache {
	cache := u.CloneFriendRelationCache()
	cache.AddKeys(u.getFollowerUserIDsKey(userID))
	return cache
}

func (u *friendRelationCacheRedis) DelSubscriberUserIDs(_ context.Context, userID string) cache.FriendRelationCache {
	cache := u.CloneFriendRelationCache()
	cache.AddKeys(u.getSubscriberUserIDsKey(userID))
	return cache
}
