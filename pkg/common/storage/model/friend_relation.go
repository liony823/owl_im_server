package model

import "time"

type FriendRelation struct {
	OwnerUserID   string    `bson:"owner_user_id"`
	RelatedUserID string    `bson:"related_user_id"`
	IsBlocked     int32     `bson:"is_blocked"`
	IsFollowing   int32     `bson:"is_following"`
	IsSubscribed  int32     `bson:"is_subscribed"`
	CreateTime    time.Time `bson:"create_time"`
	UpdateTime    time.Time `bson:"update_time"`
}
