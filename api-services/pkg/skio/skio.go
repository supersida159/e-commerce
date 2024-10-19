package skio

import (
	"net"
	"net/http"
	"net/url"

	"github.com/supersida159/e-commerce/api-services/common"
)

type Conn interface {
	ID() string
	Close() error
	URL() url.URL
	LocalAddr() net.Addr
	RemoteHeader() http.Header

	//Context of this  connection, you can save one context for one connection
	// and share it between all handlers. the handlers is called in one goroutine
	// so no need to lock context iff it only be accessed in one connection
	Context() interface{}
	SetContext(v interface{})
	Namespace() string
	Emit(event string, v ...interface{})

	//broadcast

	Join(room string)
	Leave(room string)
	LeaveAll()
	Rooms() []string
}

type AppSocket interface {
	Conn
	common.Requester
}

type appSocket struct {
	Conn
	common.Requester
}

func NewAppSocket(conn Conn, requester common.Requester) AppSocket {
	return &appSocket{
		Conn:      conn,
		Requester: requester,
	}
}
