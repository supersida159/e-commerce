package skio

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/supersida159/e-commerce/pkg/app_context"
	"github.com/supersida159/e-commerce/pkg/tokenprovider/jwt"
	"github.com/supersida159/e-commerce/src/users/repository_user"
	"github.com/supersida159/e-commerce/src/users/route_user/skuser"
)

type RealTimeEngine interface {
	UserSocket(userId int) []AppSocket
	EmitToRoom(room, event string, v ...interface{}) error
	EmitToUser(userId int, event string, v ...interface{}) error
	Run(ctx app_context.Appcontext, engine *gin.Engine) error
}

type RtEngine struct {
	server  *socketio.Server
	storage map[int][]AppSocket
	locker  *sync.RWMutex
}

// UserSocket implements RealTimeEngine.
func (e *RtEngine) UserSocket(userId int) []AppSocket {
	return e.storage[userId]
}

func NewEngine() *RtEngine {
	return &RtEngine{
		storage: make(map[int][]AppSocket),
		locker:  &sync.RWMutex{},
	}
}

func (e *RtEngine) saveAppSocket(userId int, appSck AppSocket) {
	e.locker.Lock()
	defer e.locker.Unlock()
	if v, ok := e.storage[userId]; ok {
		e.storage[userId] = append(v, appSck)
	} else {
		e.storage[userId] = []AppSocket{appSck}
	}
}

func (e *RtEngine) getAppSockets(userId int) []AppSocket {
	e.locker.RLock()
	defer e.locker.RUnlock()
	fmt.Println(e.storage)
	return e.storage[userId]
}

func (e *RtEngine) removeAppSocket(userId int, appSck AppSocket) {
	e.locker.Lock()
	defer e.locker.Unlock()
	if v, ok := e.storage[userId]; ok {
		for i := 0; i < len(v); i++ {
			if v[i] == appSck {
				e.storage[userId] = append(v[:i], v[i+1:]...)
				break
			}
		}
	}
}

func (e *RtEngine) UserSockets(userId int) []AppSocket {
	var sockets []AppSocket
	if scks, ok := e.storage[userId]; ok {
		sockets = scks
	}
	return sockets
}

func (e *RtEngine) EmitToRoom(room, event string, v ...interface{}) error {
	e.server.BroadcastToRoom("/", room, event, v...)
	return nil
}

func (e *RtEngine) EmitToUser(userId int, event string, data ...interface{}) error {
	sockets := e.getAppSockets(userId)
	for _, s := range sockets {
		s.Emit(event, data...)
	}
	return nil
}

func (e *RtEngine) Run(appctx app_context.Appcontext, engine *gin.Engine) error {
	e.server = socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			websocket.Default,
		},
	})

	e.server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID(), "IP:", s.RemoteAddr())
		return nil
	})
	e.server.OnError("/", func(s socketio.Conn, err error) {
		fmt.Println("meet error:", err)
	})

	e.server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	e.server.OnEvent("/", "Authenticate", func(s socketio.Conn, token string) {
		db := appctx.GetMainDBConnection()
		store := repository_user.NewSQLStore(db)
		tokenProvider := jwt.NewJwtProvider(appctx.GetSecretKey())
		fmt.Println("token:", token)
		payload, err := tokenProvider.Validate(token)

		if err != nil {
			s.Emit("Authenticate", err.Error())
			s.Close()
			return
		}
		user, err := store.FindUser(context.Background(), map[string]interface{}{"id": payload.UserId})

		if err != nil {
			s.Emit("Authenticate", err.Error())
			s.Close()
			return
		}
		fmt.Println("payload:", payload)
		if user.Status == 0 {
			fmt.Println("user has been deleted of baned")
			s.Emit("Authenticate", errors.New("user has been deleted of baned"))
			s.Close()
			return
		}
		appSck := NewAppSocket(s, user)
		e.saveAppSocket(user.ID, appSck)
		fmt.Println(user)
		s.Emit("Authenticate", user)
		user.Mask(false)
		s.Emit("Authenticate", user)
		e.server.OnEvent("/", "OnUserUpdateLocation", skuser.OnUserUpdateLocation(appctx, user))
	})
	go e.server.Serve()

	engine.GET("/socket.io/*any", gin.WrapH(e.server))
	engine.POST("/socket.io/*any", gin.WrapH(e.server))
	return nil
}
