package timers

/*
 * [file desc]
 * author: CC
 * email : 151503324@qq.com
 * date  : 2021.08.06
 */

import (
	"time"

	"github.com/crazycloudcc/ccgo/debugger"
)

/************************************************************************/
// constants, variables, structs, interfaces.
/************************************************************************/

/************************************************************************/
// export functions.
/************************************************************************/

// add timeout task.
func SetTimeout(duration time.Duration, callback TimedTaskCallback) *TimedTask {
	return newTimedTask(0, duration, callback)
}

// add interval task.
func SetInterval(duration time.Duration, callback TimedTaskCallback) *TimedTask {
	return newTimedTask(1, duration, callback)
}

/************************************************************************/
// moudule functions.
/************************************************************************/

/************************************************************************/
// unit tests.
/************************************************************************/

func UnitTest_TimedTask() {

	// LogDebug("SetTimeout function. ", time.Now())
	// SetTimeout(100*1e9, func() {
	// 	LogDebug("SetTimeout function. ", time.Now())
	// })

	debugger.LogDebug("SetInterval function. ", time.Now())
	SetInterval(1*1e9, func() {
		debugger.LogDebug("SetInterval function. ", time.Now())
	})
}
