package netwebsocket

/*
 * netwebsocket connection.
 * author: CC
 * email : 151503324@qq.com
 * date  : 2017.06.17
 */

import (
	"sync"

	"github.com/crazycloudcc/ccgo/debugger"

	"github.com/gorilla/websocket"
)

/************************************************************************/
// constants, variables, structs, interfaces.
/************************************************************************/

// Client TODO.
type Client struct {
	sync.RWMutex
	id      int32
	conn    *websocket.Conn
	serv    *Service
	chanRes chan []byte
}

/************************************************************************/
// export functions.
/************************************************************************/

/************************************************************************/
// moudule functions.
/************************************************************************/

// new connection.
func newClient(id int32, c *websocket.Conn, serv *Service) *Client {
	tc := new(Client)
	tc.id = id
	tc.conn = c
	tc.serv = serv
	tc.chanRes = make(chan []byte)

	go tc.read()
	go tc.write()

	return tc
}

// close
func (owner *Client) close() {
	if owner.id == -1 {
		debugger.LogDebug("client.go closed -------", owner.id, owner.chanRes)
		return
	}
	debugger.LogDebug("client.go close -------", owner.id, owner.chanRes)
	owner.conn.Close()
	close(owner.chanRes)
	owner.id = -1
}

// send data to connection.
func (owner *Client) send(buf []byte) {
	owner.Lock()
	defer owner.Unlock()
	owner.chanRes <- buf
}

func (owner *Client) read() {
	defer func() {
		debugger.LogDebug("client read unregister: ", owner.id)
		owner.serv.GetClientMngr().unregister <- owner.id
	}()

	for {
		if owner.id == -1 {
			break
		}
		_, buf, err := owner.conn.ReadMessage()
		if err != nil {
			debugger.LogDebug("client read error: ", owner.id, err)
			break
		}
		debugger.LogDebug("client read data: ", owner.id, len(buf))
		msg := bytesToMsg(buf)
		owner.serv.GetRouter().Route(owner.id, msg.ID, msg.MetaData)
	}
}

func (owner *Client) write() {
	defer func() {
		debugger.LogDebug("client write error: ", owner.id)
		owner.serv.GetClientMngr().unregister <- owner.id
	}()

	for {
		if owner.id == -1 {
			break
		}
		select {
		case buf := <-owner.chanRes:
			if owner.id == -1 || buf == nil {
				break
			}

			// debugger.LogDebug("client write data: ", owner.id, buf)
			owner.conn.WriteMessage(websocket.BinaryMessage, buf)
		}
	}
}

/************************************************************************/
// unit tests.
/************************************************************************/
