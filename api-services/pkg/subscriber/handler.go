package subscriber

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"strings"

// 	"github.com/supersida159/e-commerce/api-services/common"
// 	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
// 	"github.com/supersida159/e-commerce/api-services/pkg/kafka/consumerlocal"
// 	"github.com/supersida159/e-commerce/api-services/pkg/pubsub"
// )

// func RunExpiredUpdateOrder(appCtx app_context.AppContext) consumerJob {
// 	return consumerJob{
// 		Title: "RunExpiredUpdateOrder",
// 		Hld: func(ctx context.Context, msg *pubsub.Message) *common.AppError {
// 			sagaID := msg.Data().(string)

// 			sagaID = strings.TrimPrefix(sagaID, "Saga-")

// 			sagaConsumer := appCtx.GetConsumer()

// 			updateChannel, err := sagaConsumer.GetEventChannel(sagaID, consumerlocal.UpdateChannel)
// 			if err != nil {
// 				return common.ErrInternalServerError(fmt.Errorf("failed to get update channel for saga %s: %w", sagaID, err))
// 			}

// 			select {
// 			case event := <-updateChannel:
// 				log.Printf("Received event for expired saga %s: %v", sagaID, event)
// 				// Add any necessary processing logic
// 			case <-ctx.Done():
// 				return common.ErrInternalServerError(ctx.Err())
// 			default:
// 				log.Printf("No event found for saga %s", sagaID)
// 			}

// 			return nil
// 		},
// 	}
// }
