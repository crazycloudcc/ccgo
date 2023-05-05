package nettcp

/*
 * network tcp connection.
 * author: CC
 * email : 151503324@qq.com
 * date  : 2017.06.17
 */

import (
	"bytes"
	"ccgo/datastructs"
	"ccgo/debugger"
	"fmt"
	"net"
	"sync"
	"time"
)

/************************************************************************/
// constants, variables, structs, interfaces.
/************************************************************************/

// MsgChanMax message channel size.
const MsgChanMax int32 = 100

// Client struct.
type Client struct {
	sync.Mutex
	id            int32             // Client id.
	conn          net.Conn          // connection instance.
	buffer        *bytes.Buffer     // recv buffer.
	chanMsg       *datastructs.Chan // message channel.
	readLen       int32             // [analysist] all read data length.
	writeLen      int32             // [analysist] all write data length.
	serv          *Service          // network Service instance.
	heartbeatTime time.Time         // last recv heart beat message time.
}

/************************************************************************/
// export functions.
/************************************************************************/

// ConnID get connection id.
func (owner *Client) ConnID() int32 {
	return owner.id
}

// Run ready to recv data.
func (owner *Client) Run(id int32) {
	owner.id = id
	go owner.doRead()
	go owner.heartbeat()
}

// Close connection by external function.
func (owner *Client) Close() {
	if owner.id == -1 {
		return
	}
	debugger.LogDebug("client.go - Close: ", owner.id)
	owner.serv.GetClientMngr().DelUID(owner.id)
}

// Send data to connection.
func (owner *Client) Send(msg interface{}) {
	m := msg.(*Msg)
	debugger.LogDebug(fmt.Sprintf("client.go - send msglen:[%v], msgid:[%v]", m.Len, m.ID))

	buf, err := Marshal(m)
	if err != nil {
		debugger.LogError("client.go - send error: ", owner.id, err)
		return
	}
	owner.doWrite(buf)
}

/************************************************************************/
// moudule functions.
/************************************************************************/

// new connection.
func newClient(c net.Conn, serv *Service) *Client {
	tc := new(Client)
	tc.id = -1
	tc.conn = c
	tc.buffer = new(bytes.Buffer)
	tc.chanMsg = datastructs.NewChan(MsgChanMax, tc.doParse)
	tc.readLen = 0
	tc.writeLen = 0
	tc.serv = serv
	return tc
}

// close connection by connection error.
func (owner *Client) doClose() {
	if owner.id == -1 {
		return
	}
	owner.id = -1
	owner.conn.Close()
	owner.chanMsg.Close()
	debugger.LogDebug("client.go - doClose: ", owner.id)
}

// goroutine recv data from connection.
func (owner *Client) doRead() {
	for {
		buf := make([]byte, BufReadMax)
		len, err := owner.conn.Read(buf)
		if err != nil {
			if err.Error() != "EOF" {
				debugger.LogWarn("client.go - read error: ", owner.id, err)
			}
			break
		}
		owner.readLen += int32(len)
		owner.chanMsg.Write(buf[:len])
		debugger.LogDebug(fmt.Sprintf("client.go - id: [%d] read len: [%d], total read len: [%d]", owner.id, len, owner.readLen))
	}
	owner.Close()
}

// send data to connection.
func (owner *Client) doWrite(buf []byte) {
	owner.Lock()
	defer owner.Unlock()
	len, err := owner.conn.Write(buf)
	if err != nil {
		debugger.LogError("client.go - write error: ", owner.id, err)
		owner.Close()
		return
	}
	owner.writeLen += int32(len)
	debugger.LogDebug(fmt.Sprintf("client.go - id: [%d] write len: [%d], total write len: [%d]", owner.id, len, owner.writeLen))
}

// goroutine heart beat.
func (owner *Client) heartbeat() {
	for {
		time.Sleep(time.Duration(HeartBeatTimeTick) * 1e9)
		d := time.Since(owner.heartbeatTime)
		if d > time.Duration(HeartBeatTimeTick) {
			break
		}
	}
	owner.Close()
}

// refresh heartbeat time.
func (owner *Client) refreshHBTime() {
	owner.heartbeatTime = time.Now()
}

// parse and route connection data.
func (owner *Client) doParse(b interface{}) {
	owner.buffer.Write(b.([]byte))
	for {
		if owner.buffer.Len() < MsgHeadLen {
			// debugger.LogWarn(fmt.Sprintf("client.go - owner.buffer.Len()[%d] < MsgHeadLen[%d]", owner.buffer.Len(), MsgHeadLen))
			break
		}
		msgLen := owner.buffer.Bytes()[:MsgLenIndexE]
		l, err := ReadMsgLen(msgLen)
		if err != nil {
			debugger.LogError("client.go - ReadMsgLen error: ", owner.id, err)
			owner.Close()
			break
		}
		if int(l)+MsgHeadLen > owner.buffer.Len() {
			// debugger.LogWarn(fmt.Sprintf("client.go - int(l)+MsgHeadLen [%d] > owner.buffer.Len() [%d]", int(l)+MsgHeadLen, owner.buffer.Len()))
			break
		}

		buf := make([]byte, int(l)+MsgHeadLen)
		owner.buffer.Read(buf)
		msg, err := Unmarshal(buf)
		if err != nil {
			debugger.LogError("client.go - Unmarshal error: ", owner.id, err)
			owner.Close()
			break
		}

		owner.serv.GetRouter().Route(owner.id, msg.ID, msg.MetaData)
	}
}

/************************************************************************/
// unit tests.
/************************************************************************/
