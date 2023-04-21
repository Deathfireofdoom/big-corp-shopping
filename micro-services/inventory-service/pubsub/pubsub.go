package pubsub

import(
	"context"
	"inventory/entity"
)

type PubSub interface {
	Consume(ctx context.Context, ch chan<- entity.InventoryRequest, failOnError bool) error
}






// TODO add this
//type ConsumerSQS struct {
//
//}