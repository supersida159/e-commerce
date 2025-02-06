package gin_order

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/pkg/kafka/saga"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WebSocketOrderStatusHandler(appCtx app_context.AppContext, saga *saga.Orchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := r.URL.Query().Get("orderId")
		if orderID == "" {
			http.Error(w, "Missing orderId", http.StatusBadRequest)
			return
		}

		// Upgrade the HTTP connection to a WebSocket connection
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		// Use the existing orchestrator to subscribe to order updates
		saga.SubscribeToOrderUpdates(orderID, conn)
	}
}
