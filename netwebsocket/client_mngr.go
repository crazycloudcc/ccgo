package netwebsocket

import (
	"sync"

	"github.com/crazycloudcc/ccgo/debugger"

	"github.com/gorilla/websocket"
)

/*
 * 网络连接管理.
 * 管理Conn类型实例.
 */

/************************************************************************/
// constants, variables, structs, interfaces.
/************************************************************************/

// NetWSCallback TODO
type NetWSCallback func(connID int32)

// ClientMngr TODO.
type ClientMngr struct {
	sync.RWMutex
	serv               *Service
	seed               int32
	clients            map[int32]*Client
	broadcast          chan []byte
	register           chan *websocket.Conn
	unregister         chan int32
	callbackDisConnect NetWSCallback
	callbackConnect    NetWSCallback
}

/************************************************************************/
// export functions.
/************************************************************************/

// SendToConn 通过连接ID发送
func (owner *ClientMngr) SendToConn(ConnID int32, msg *Msg) {
	buf := msgToBytes(msg)
	for id, client := range owner.clients {
		if id == ConnID {
			client.send(buf)
			break
		}
	}
}

// SendToAll TODO
func (owner *ClientMngr) SendToAll(msg *Msg) {
	buf := msgToBytes(msg)
	for _, client := range owner.clients {
		client.send(buf)
	}
}

// SendWithOut TODO
func (owner *ClientMngr) SendWithOut(msg *Msg, ignore int32) {
	buf := msgToBytes(msg)
	for id, client := range owner.clients {
		if id != ignore {
			client.send(buf)
		}
	}
}

// DisConnect TODO
func (owner *ClientMngr) DisConnect(connID int32) {
	for id, client := range owner.clients {
		if id != connID {
			client.close()
			break
		}
	}
}

// SetConnectCallback TODO
func (owner *ClientMngr) SetConnectCallback(callback NetWSCallback) {
	owner.callbackConnect = callback
}

// SetDisConnectCallback TODO
func (owner *ClientMngr) SetDisConnectCallback(callback NetWSCallback) {
	owner.callbackDisConnect = callback
}

/************************************************************************/
// moudule functions.
/************************************************************************/

// create new ClientMngr.
func newClientMngr(serv *Service) *ClientMngr {
	tcm := new(ClientMngr)
	tcm.serv = serv
	tcm.clients = make(map[int32]*Client)
	tcm.broadcast = make(chan []byte)
	tcm.register = make(chan *websocket.Conn)
	tcm.unregister = make(chan int32)
	return tcm
}

func (owner *ClientMngr) start() {
	for {
		select {
		case conn := <-owner.register:
			owner.seedTick()
			id := owner.seed
			debugger.LogDebug("A new socket has connected. ", id, conn)
			owner.clients[id] = newClient(id, conn, owner.serv)
			owner.callbackConnect(id)
		case cid := <-owner.unregister:
			debugger.LogDebug("A socket has disconnected. ", cid)
			if client, ok := owner.clients[cid]; ok {
				debugger.LogDebug("A socket has disconnected callback ", cid)
				client.close()
				delete(owner.clients, cid)
				owner.callbackDisConnect(cid)
			}
		case message := <-owner.broadcast:
			debugger.LogDebug("broadcast message: ", message)
			// for id, client := range owner.clients {
			// 	select {
			// 	case conn.send <- message:
			// 	default:
			// 		close(conn.send)
			// 		delete(owner.clients, id)
			// 	}
			// }
		}
	}
}

// seed tick.
func (owner *ClientMngr) seedTick() {
	owner.seed++
}

/************************************************************************/
// unit tests.
/************************************************************************/
