package common

import (
	"github.com/supersida159/e-commerce/pkg/pubsub"
)

const (
	TopicUserLikeRestaurant   pubsub.Topic = "TopicUserLikeRestaurant"
	TopicUserUnLikeRestaurant pubsub.Topic = "TopicUserUnLikeRestaurant"
	TopicUploadImg            pubsub.Topic = "TopicUploadImg"
)
