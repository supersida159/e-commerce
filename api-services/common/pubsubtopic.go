package common

import "github.com/supersida159/e-commerce/api-services/pkg/pubsub"

const (
	TopicUserLikeRestaurant   pubsub.Topic = "TopicUserLikeRestaurant"
	TopicUserUnLikeRestaurant pubsub.Topic = "TopicUserUnLikeRestaurant"
	TopicUploadImg            pubsub.Topic = "TopicUploadImg"
	TopicOrderCreated         pubsub.Topic = "TopicOrderCreated"
	TopicOrderCancelled       pubsub.Topic = "TopicOrderCancelled"
	TopicOrderUpdated         pubsub.Topic = "TopicOrderUpdated"
	TopicOrderDeleted         pubsub.Topic = "TopicOrderDeleted"
	TopicOrderConfirmed       pubsub.Topic = "TopicOrderConfirmed"
	TopicOrderExpired         pubsub.Topic = "TopicOrderExpired"
)
