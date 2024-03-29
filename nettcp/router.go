package nettcp

/*
 * message router.
 * author: CC
 * email : 151503324@qq.com
 * date  : 2017.06.17
 */

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/crazycloudcc/ccgo/datastructs"
	"github.com/crazycloudcc/ccgo/debugger"
	"github.com/crazycloudcc/ccgo/timers"
)

/************************************************************************/
// constants, variables, structs, interfaces.
/************************************************************************/

// IHandler interface.
type IHandler interface {
	Handle(ev *Event)
}

// Event struct.
type Event struct {
	Connid   int32
	Msgid    uint32
	MetaData []byte
}

// Router struct.
type Router struct {
	sync.Mutex
	msgHandlers  *datastructs.Hash
	chanMsg      *datastructs.Chan
	timeRecorder *timers.TimeRecorder
}

/************************************************************************/
// export functions.
/************************************************************************/

// RegHandler 注册消息处理函数.
// 参数: 1-消息ID, 2-消息处理器结构体实例.
func (owner *Router) RegHandler(msgid uint32, handler interface{}) {
	debugger.LogDebug("router.go - register hander with msgid: ", msgid)
	owner.msgHandlers.Add(msgid, reflect.TypeOf(handler))
}

// Route 解析后的消息放入路由中的事件队列进行分发.
func (owner *Router) Route(connid int32, msgid uint32, metadata []byte) {
	owner.chanMsg.Write(&Event{connid, msgid, metadata})
}

/************************************************************************/
// moudule functions.
/************************************************************************/

// 创建消息路由器.
func newRouter(max int32) *Router {
	r := new(Router)
	r.msgHandlers = datastructs.NewHash(max)
	r.chanMsg = datastructs.NewChan(max, r.doRoute)
	r.timeRecorder = timers.NewTimeRecorder()
	return r
}

// 事件分发.
func (owner *Router) doRoute(data interface{}) {
	ev := data.(*Event)
	debugger.LogDebug(fmt.Sprintf("route.go - connid:[%d] msgid:[%d] len:[%d]", ev.Connid, ev.Msgid, len(ev.MetaData)))
	handler := owner.msgHandlers.Get(ev.Msgid)

	if handler != nil {
		////////////////////////////////////////////////////////
		// 调试.
		// s := time.Now()
		h := reflect.New(handler.(reflect.Type)).Interface()
		h.(IHandler).Handle(ev)
		// d := time.Now().Sub(s)
		// debugger.LogDebug(fmt.Sprintf("router.go - doRoute connid:[%d] msgid:[%d] cost time:[%v]", ev.Connid, ev.Msgid, d))
		////////////////////////////////////////////////////////
		return
	}
	debugger.LogWarn("router.go - not found msg handler: ", ev.Connid, ev.Msgid)
}

/************************************************************************/
// unit tests.
/************************************************************************/
