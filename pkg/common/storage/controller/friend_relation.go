package controller

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

type FriendRelationDatabase interface {
	Create(ctx context.Context, relations []*model.FriendRelation) error
	Delete(ctx context.Context, relations []*model.FriendRelation) error
	BlockUser(ctx context.Context, userID string, blockUserID string) error
	UnblockUser(ctx context.Context, userID string, blockUserID string) error
	FollowUser(ctx context.Context, userID string, relatedUserID string) error
	UnfollowUser(ctx context.Context, userID string, relatedUserID string) error
	SubscribeUser(ctx context.Context, userID string, relatedUserID string) error
	UnsubscribeUser(ctx context.Context, userID string, relatedUserID string) error
	GetFollowerUserIDs(ctx context.Context, userID string) ([]string, error)
	GetSubscriberUserIDs(ctx context.Context, userID string) ([]string, error)
}

type friendRelationDatabase struct {
	friendRelation database.FriendRelation
	cache          cache.FriendRelationCache
}

func NewFriendRelation(friendRelation database.FriendRelation, cache cache.FriendRelationCache) FriendRelationDatabase {
	return &friendRelationDatabase{friendRelation: friendRelation, cache: cache}
}

func (r *friendRelationDatabase) Create(ctx context.Context, relations []*model.FriendRelation) error {
	if err := r.friendRelation.Create(ctx, relations); err != nil {
		return err
	}
	return r.deleteRelatedUserIDsCache(ctx, relations)
}

func (r *friendRelationDatabase) Delete(ctx context.Context, relations []*model.FriendRelation) error {
	if err := r.friendRelation.Delete(ctx, relations); err != nil {
		return err
	}
	return r.deleteRelatedUserIDsCache(ctx, relations)
}

func (r *friendRelationDatabase) deleteRelatedUserIDsCache(ctx context.Context, relations []*model.FriendRelation) error {
	cache := r.cache.CloneFriendRelationCache()
	for _, relation := range relations {
		cache = cache.DelFollowerUserIDs(ctx, relation.OwnerUserID)
		cache = cache.DelSubscriberUserIDs(ctx, relation.OwnerUserID)
	}
	return cache.ChainExecDel(ctx)
}

func (r *friendRelationDatabase) BlockUser(ctx context.Context, userID string, blockUserID string) error {
	if err := r.friendRelation.UpdateByMap(ctx, userID, blockUserID, map[string]any{
		"is_blocked": 1,
	}); err != nil {
		return err
	}
	return r.deleteRelatedUserIDsCache(ctx, []*model.FriendRelation{
		{OwnerUserID: userID, RelatedUserID: blockUserID},
	})
}

func (r *friendRelationDatabase) UnblockUser(ctx context.Context, userID string, blockUserID string) error {
	if err := r.friendRelation.UpdateByMap(ctx, userID, blockUserID, map[string]any{
		"is_blocked": 0,
	}); err != nil {
		return err
	}
	return r.deleteRelatedUserIDsCache(ctx, []*model.FriendRelation{
		{OwnerUserID: userID, RelatedUserID: blockUserID},
	})
}

func (r *friendRelationDatabase) FollowUser(ctx context.Context, userID string, relatedUserID string) error {
	if err := r.friendRelation.UpdateByMap(ctx, userID, relatedUserID, map[string]any{
		"is_followed": 1,
	}); err != nil {
		return err
	}
	return r.deleteRelatedUserIDsCache(ctx, []*model.FriendRelation{
		{OwnerUserID: userID, RelatedUserID: relatedUserID},
	})
}

func (r *friendRelationDatabase) UnfollowUser(ctx context.Context, userID string, relatedUserID string) error {
	if err := r.friendRelation.UpdateByMap(ctx, userID, relatedUserID, map[string]any{
		"is_followed": 0,
	}); err != nil {
		return err
	}
	return r.deleteRelatedUserIDsCache(ctx, []*model.FriendRelation{
		{OwnerUserID: userID, RelatedUserID: relatedUserID},
	})
}

func (r *friendRelationDatabase) SubscribeUser(ctx context.Context, userID string, relatedUserID string) error {
	if err := r.friendRelation.UpdateByMap(ctx, userID, relatedUserID, map[string]any{
		"is_subscribed": 1,
	}); err != nil {
		return err
	}
	return r.deleteRelatedUserIDsCache(ctx, []*model.FriendRelation{
		{OwnerUserID: userID, RelatedUserID: relatedUserID},
	})
}

func (r *friendRelationDatabase) UnsubscribeUser(ctx context.Context, userID string, relatedUserID string) error {
	if err := r.friendRelation.UpdateByMap(ctx, userID, relatedUserID, map[string]any{
		"is_subscribed": 0,
	}); err != nil {
		return err
	}
	return r.deleteRelatedUserIDsCache(ctx, []*model.FriendRelation{
		{OwnerUserID: userID, RelatedUserID: relatedUserID},
	})
}

func (r *friendRelationDatabase) GetFollowerUserIDs(ctx context.Context, userID string) ([]string, error) {
	return r.cache.GetFollowerUserIDs(ctx, userID)
}

func (r *friendRelationDatabase) GetSubscriberUserIDs(ctx context.Context, userID string) ([]string, error) {
	return r.cache.GetSubscriberUserIDs(ctx, userID)
}
