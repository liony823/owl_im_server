package mgo

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/mongoutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FriendRelationMgo struct {
	coll *mongo.Collection
}

func (u *FriendRelationMgo) userRelationFilter(ownerUserID, relatedUserID string) bson.M {
	return bson.M{
		"owner_user_id":   ownerUserID,
		"related_user_id": relatedUserID,
	}
}

func (u *FriendRelationMgo) userRelationsFilter(relations []*model.FriendRelation) bson.M {
	if len(relations) == 0 {
		return nil
	}
	or := make(bson.A, 0, len(relations))
	for _, relation := range relations {
		or = append(or, u.userRelationFilter(relation.OwnerUserID, relation.RelatedUserID))
	}
	return bson.M{"$or": or}
}

func NewFriendRelationMgo(db *mongo.Database) (database.FriendRelation, error) {
	coll := db.Collection(database.FriendRelationName)
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "owner_user_id", Value: 1},
			{Key: "related_user_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}
	return &FriendRelationMgo{coll: coll}, nil
}

func (u *FriendRelationMgo) Create(ctx context.Context, relations []*model.FriendRelation) error {
	for i, relation := range relations {
		if relation.CreateTime.IsZero() {
			relations[i].CreateTime = time.Now()
		}
		if relation.UpdateTime.IsZero() {
			relations[i].UpdateTime = time.Now()
		}
	}

	return mongoutil.InsertMany(ctx, u.coll, relations)
}

func (u *FriendRelationMgo) Take(ctx context.Context, userID, relatedUserID string) (*model.FriendRelation, error) {
	filter := bson.M{
		"owner_user_id":   userID,
		"related_user_id": relatedUserID,
	}
	return mongoutil.FindOne[*model.FriendRelation](ctx, u.coll, filter)
}

func (u *FriendRelationMgo) UpdateByMap(ctx context.Context, userID, relatedUserID string, args map[string]any) error {
	filter := bson.M{
		"owner_user_id":   userID,
		"related_user_id": relatedUserID,
	}
	if args["update_time"] == nil {
		args["update_time"] = time.Now()
	}
	return mongoutil.UpdateOne(ctx, u.coll, filter, bson.M{"$set": args}, false)
}

func (u *FriendRelationMgo) Delete(ctx context.Context, relations []*model.FriendRelation) error {
	if len(relations) == 0 {
		return nil
	}
	return mongoutil.DeleteMany(ctx, u.coll, u.userRelationsFilter(relations))
}

func (u *FriendRelationMgo) GetFollowerUserIDs(ctx context.Context, userID string) ([]string, error) {
	filter := bson.M{
		"owner_user_id": userID,
		"is_following":  1,
		"is_blocked":    0,
	}
	return mongoutil.Find[string](ctx, u.coll, filter, options.Find().SetProjection(bson.M{"_id": 0, "related_user_id": 1}))
}

func (u *FriendRelationMgo) GetSubscriberUserIDs(ctx context.Context, userID string) ([]string, error) {
	filter := bson.M{
		"owner_user_id": userID,
		"is_subscribed": 1,
		"is_blocked":    0,
	}
	return mongoutil.Find[string](ctx, u.coll, filter, options.Find().SetProjection(bson.M{"_id": 0, "related_user_id": 1}))
}
