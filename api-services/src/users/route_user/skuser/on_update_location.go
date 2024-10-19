package skuser

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
)

type LocationData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func OnUserUpdateLocation(appCtx app_context.Appcontext, requester common.Requester) func(s socketio.Conn, location LocationData) {
	return func(s socketio.Conn, location LocationData) {
		log.Println("User:", s.ID(), "UpdateLocation: ", location)
		log.Println("User requester:", requester)
	}
}
