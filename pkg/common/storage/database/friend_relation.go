package database

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

type FriendRelation interface {
	Create(ctx context.Context, relations []*model.FriendRelation) error
	Delete(ctx context.Context, relations []*model.FriendRelation) error
	Take(ctx context.Context, userID, related_user_id string) (*model.FriendRelation, error)
	UpdateByMap(ctx context.Context, userID, related_user_id string, args map[string]any) error
	GetFollowerUserIDs(ctx context.Context, userID string) ([]string, error)
	GetSubscriberUserIDs(ctx context.Context, userID string) ([]string, error)
}
